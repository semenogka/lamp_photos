import torch
import clip
from PIL import Image
import json
import streamlit as st
import numpy as np

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



text_tokens = clip.tokenize(labels).to(device)
with torch.no_grad():
    text_features = model.encode_text(text_tokens)
    text_features /= text_features.norm(dim=-1, keepdim=True)

def get_image_features_and_label(path):
    embs = []
    best_labels = []

    for angle in range(0, 360, 90):
        file = Image.open(path).convert("RGB").rotate(angle, expand=True)
        mirrored = file.transpose(Image.FLIP_LEFT_RIGHT)

        imageOrig = preprocess(file).unsqueeze(0).to(device)
        imageMir = preprocess(mirrored).unsqueeze(0).to(device)
        with torch.no_grad():
            image_features_orig = model.encode_image(imageOrig)
            image_features_orig /= image_features_orig.norm(dim=-1, keepdim=True)
            similarities_orig = (image_features_orig @ text_features.T).squeeze(0)
            best_idx_orig = similarities_orig.argmax().item()
            best_label_orig = labels[best_idx_orig]
            text_orig = clip.tokenize([best_label_orig]).to(device)
            tf_orig = model.encode_text(text_orig)
            tf_orig /= tf_orig.norm(dim=-1, keepdim=True)

            image_features_mir = model.encode_image(imageMir)
            image_features_mir /= image_features_mir.norm(dim=-1, keepdim=True)
            similarities_mir = (image_features_mir @ text_features.T).squeeze(0)
            best_idx_mir = similarities_mir.argmax().item()
            best_label_mir = labels[best_idx_mir]
            text_mir = clip.tokenize([best_label_mir]).to(device)
            tf_mir = model.encode_text(text_mir)
            tf_mir /= tf_mir.norm(dim=-1, keepdim=True)

            embs.append(image_features_orig)
            embs.append(image_features_mir)
            best_labels.append(tf_orig)
            best_labels.append(tf_mir)
    return embs, best_labels
            



result = []

embs, label = get_image_features_and_label("image.png")

with open(fr"alldata\ndata\npodves.json", "r", encoding="utf-8") as f:
    data = json.load(f)

for item in data:
    minRes = []

    text2 = clip.tokenize([item["label"]]).to(device)
    with torch.no_grad():
        tf2 = model.encode_text(text2)
        tf2 /= tf2.norm(dim=-1, keepdim=True)

    emb = torch.tensor(item['emb']).to(device)

    for j in range(len(embs)):
        with torch.no_grad():
            tf1 = label[j]

            text_sim = torch.cosine_similarity(tf1, tf2).item()

            if text_sim > 0.97:
                emb_user = embs[j]
                image_sim = torch.cosine_similarity(emb_user, emb, dim=-1).item()
                minRes.append(image_sim)

    if minRes:
        good_sim = max(minRes)
        if good_sim > 0.7:
            result.append({
                "sim": good_sim,
                "link": item['link'],
            })
            

result_sorted = sorted(result, key=lambda x: x["sim"], reverse=True)
for r in range(60):
    print(result_sorted[r])
