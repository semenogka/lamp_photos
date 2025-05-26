import os
import torch
import clip
from PIL import Image
import json
import streamlit as st
import numpy as np
import requests
from io import BytesIO

# Устройство и загрузка модели
device = "cuda" if torch.cuda.is_available() else "cpu"
model, preprocess = clip.load("ViT-B/32", device)

st.title("Поиск похожих светильников")

# Загрузка изображений
uploaded_files = st.file_uploader(
    "Загрузите или вставьте до 3 изображений",
    type=["png", "jpg", "jpeg"],
    accept_multiple_files=True
)

# Список текстовых меток
labels = [
    "светильник в форме круга", "светильник с круглой формой корпуса", "плоский круглый потолочный светильник",
    "прямоугольный светильник", "светильник вытянутой прямоугольной формы", "тонкий прямоугольный потолочный светильник",
    "квадратный светильник", "светильник с квадратным корпусом", "компактный квадратный светильник",
    "светильник в форме шара", "шарообразный подвесной светильник", "объемный светильник, похожий на шар",
    "светильник в форме куба", "кубообразный корпус светильника", "геометрически строгий светильник в форме куба",
    "длинный продолговатый светильник", "узкий светильник в форме полоски", "линейный светильник вытянутой формы",
    "светильник с угловатой формой", "светильник с чёткими углами и гранями", "светильник с резкой геометрией корпуса",
    "конусообразный абажур", "светильник с колоколообразным плафоном", "абажур, сужающийся кверху",
    "цилиндрический светильник", "вертикальный цилиндр в виде светильника", "корпус светильника в форме трубы"
]

text_tokens = clip.tokenize(labels).to(device)
with torch.no_grad():
    text_features = model.encode_text(text_tokens)
    text_features /= text_features.norm(dim=-1, keepdim=True)

def get_image_features_and_label(image: Image.Image):
    embs, best_text_features = [], []

    for angle in range(0, 360, 90):
        file = image.rotate(angle, expand=True)
        mirrored = file.transpose(Image.FLIP_LEFT_RIGHT)

        imageOrig = preprocess(file).unsqueeze(0).to(device)
        imageMir = preprocess(mirrored).unsqueeze(0).to(device)

        with torch.no_grad():
            for img_tensor in [imageOrig, imageMir]:
                image_features = model.encode_image(img_tensor)
                image_features /= image_features.norm(dim=-1, keepdim=True)
                similarities = (image_features @ text_features.T).squeeze(0)
                best_idx = similarities.argmax().item()
                best_text_features.append(text_features[best_idx].unsqueeze(0))
                embs.append(image_features)

    return embs, best_text_features

# Категории и пути к JSON
base_path = os.path.join("alldata", "ndata")

buttons_files = {
    "Акссесуары": os.path.join(base_path, "naks.json"),
    "Архитектурные": os.path.join(base_path, "narchitect.json"),
    "Бра и настенные": os.path.join(base_path, "nbra.json"),
    "Фасадные": os.path.join(base_path, "nfasad.json"),
    "Грунтовые": os.path.join(base_path, "ngroont.json"),
    "Ландшафтные": os.path.join(base_path, "nland.json"),
    "Люстры": os.path.join(base_path, "nlystri.json"),
    "Магнитные трековые": os.path.join(base_path, "nmagn.json"),
    "Настольные": os.path.join(base_path, "nnastol.json"),
    "Настенно-потолочные": os.path.join(base_path, "nnastpot.json"),
    "Парковые": os.path.join(base_path, "nparkovie.json"),
    "Подсветные": os.path.join(base_path, "npodsvet.json"),
    "Подвесные": os.path.join(base_path, "npodves.json"),
    "Потолочные": os.path.join(base_path, "npotol.json"),
    "Прожекторы": os.path.join(base_path, "nprojectors.json"),
    "Точечные накладные": os.path.join(base_path, "ntochnakl.json"),
    "Точечные подвесные": os.path.join(base_path, "ntochpodv.json"),
    "Точечные Встроенные": os.path.join(base_path, "ntochvstr.json"),
    "Торшеры": os.path.join(base_path, "ntorsher.json"),
    "Трековые": os.path.join(base_path, "ntrack.json"),
    "Тротуарные": os.path.join(base_path, "ntrotuarnie.json"),
    "Встраиваемые": os.path.join(base_path, "nvstraivaem.json")
}

# Выбор категории
selected_categories = st.multiselect("выберите до двух категорий", list(buttons_files.keys()) ,max_selections=2)
start_search = st.button("найти похожие товары")


# Основная логика поиска
if uploaded_files and selected_categories and start_search:
    result = []
    all_embs, all_labels = [], []
    with st.spinner("Ищем похожие..."):
        for file in uploaded_files[:3]:
            image = Image.open(file).convert("RGB")
            embs, labels = get_image_features_and_label(image)
            all_embs.append(embs)
            all_labels.append(labels)

        total_items = 0 
        

        for category in selected_categories:
            
            

            clicked_file = buttons_files[category]
            with open(clicked_file, "r", encoding="utf-8") as f:
                data = json.load(f)

            for item in data:
                minRes = []
                text2 = clip.tokenize([item["label"]]).to(device)
                with torch.no_grad():
                    tf2 = model.encode_text(text2)
                    tf2 /= tf2.norm(dim=-1, keepdim=True)
                emb = torch.tensor(item['emb'], dtype=torch.float32).to(device)

                for embs_img, labels_img in zip(all_embs, all_labels):
                    for tf1, emb_user in zip(labels_img, embs_img):
                        text_sim = torch.cosine_similarity(tf1, tf2).item()
                        if text_sim > 0.95:
                            image_sim = torch.cosine_similarity(emb_user, emb, dim=-1).item()
                            minRes.append(image_sim)

                if minRes:
                    good_sim = max(minRes)
                    if good_sim > 0.7:
                        result.append({
                            "sim": good_sim,
                            "link": item['link'],
                            "name": item['name']
                        })

    result_sorted = sorted(result, key=lambda x: x["sim"], reverse=True)
    st.subheader("Найденные похожие товары:")   

    for r in result_sorted[:100]:
        cols = st.columns([1, 4])
        path = os.path.join("allimgs", r.get('name', ''))
        try:
            url = f"https://raw.githubusercontent.com/semenogka/lamp_photos/main/allimgs/{r['name']}"
            response = requests.get(url)
            img = Image.open(BytesIO(response.content))
        except Exception:
            img = None
        with cols[0]:
            if img:
                st.image(img, width=120)
            else:
                st.write("Изображение не найдено")
        with cols[1]:
            st.markdown(f"[{r['link']}]({r['link']}) — Сходство: {r['sim'] * 100:.2f}%")
