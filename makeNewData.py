import torch
import clip
import os
import json
from PIL import Image
from tqdm import tqdm 

device = "cuda" if torch.cuda.is_available() else "cpu"
model, preprocess = clip.load("ViT-B/32", device)

labels = [
    "светильник в форме круга",
    "светильник с круглой формой корпуса",
    "плоский круглый потолочный светильник",

    "прямоугольный светильник",
    "светильник вытянутой прямоугольной формы",
    "тонкий прямоугольный потолочный светильник",

    "квадратный светильник",
    "светильник с квадратным корпусом",
    "компактный квадратный светильник",

    "светильник в форме шара",
    "шарообразный подвесной светильник",
    "объемный светильник, похожий на шар",

    "светильник в форме куба",
    "кубообразный корпус светильника",
    "геометрически строгий светильник в форме куба",

    "длинный продолговатый светильник",
    "узкий светильник в форме полоски",
    "линейный светильник вытянутой формы",

    "светильник с угловатой формой",
    "светильник с чёткими углами и гранями",
    "светильник с резкой геометрией корпуса",

    "конусообразный абажур",
    "светильник с колоколообразным плафоном",
    "абажур, сужающийся кверху",

    "цилиндрический светильник",
    "вертикальный цилиндр в виде светильника",
    "корпус светильника в форме трубы"
]

files = [
    # "aks",
    # "architect",
    # "bra",
    # "fasad",
    # "groont",
    # "land",
    # "lystri",
    "magn",
    "nastol",
    "nastpot",
    "parkovie",
    "podsvet",
    "podves",
    "potol",
    "projectors",
    "tochnakl",
    "tochpodv",
    "tochvstr",
    "torsher",
    "track",
    "trotuarnie",
    "vstraivaem"
]


text_tokens = clip.tokenize(labels).to(device)
with torch.no_grad():
    textF = model.encode_text(text_tokens)
    textF /= textF.norm(dim=-1, keepdim=True)

def get_emb_n_text(path):
    try:
        image = Image.open(path).convert("RGB")
    except Exception as e:
        print(f"Ошибка при открытии файла {path}: ")
        return None, None
    image = preprocess(image).unsqueeze(0).to(device)
    with torch.no_grad():
        imageF = model.encode_image(image)
        imageF /= imageF.norm(dim=-1, keepdim=True)
        sim = (imageF @ textF.T).squeeze(0)
        best_idx = sim.argmax().item()
        best_label = labels[best_idx]
    return imageF, best_label

folder = "allimgs"
for filename in files:
    results = []
    
    with open(fr"alldata/data/{filename}.json", "r", encoding="utf-8") as f:
        data = json.load(f)

    for item in data:
        path = os.path.join(folder, item["name"])

        if not os.path.isfile(path) or not path.lower().endswith((".png", ".jpg", ".jpeg", ".webp", ".bmp")):
            print(f"SKIP: {path}")
            continue

        emb, label = get_emb_n_text(path)
        if emb is None:
            continue 
        emb = emb.cpu().tolist()
        
        results.append({
            "emb": emb,
            "label": label,
            "link": item["link"],
            "name": item["name"]
        })
    with open(fr"alldata\ndata/n{filename}.json", "w", encoding="utf-8") as x:
        json.dump(results, x, ensure_ascii=False, indent=2)
    print(filename, " ready")
