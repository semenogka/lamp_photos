import torch
import clip
from PIL import Image

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
    image = preprocess(Image.open(path).convert("RGB")).unsqueeze(0).to(device)
    with torch.no_grad():
        image_features = model.encode_image(image)
        image_features /= image_features.norm(dim=-1, keepdim=True)
        similarities = (image_features @ text_features.T).squeeze(0)
        best_idx = similarities.argmax().item()
        best_label = labels[best_idx]
        best_score = similarities[best_idx].item()
    return image_features, best_label, best_score


emb1, label1, score1 = get_image_features_and_label("image1.png")
print(f"Картинка 1: '{label1}' ({score1 * 100:.2f}%)")


emb2, label2, score2 = get_image_features_and_label("image.png")
print(f"Картинка 2: '{label2}' ({score2 * 100:.2f}%)")


text1 = clip.tokenize([label1]).to(device)
text2 = clip.tokenize([label2]).to(device)
with torch.no_grad():
    tf1 = model.encode_text(text1)
    tf2 = model.encode_text(text2)
    tf1 /= tf1.norm(dim=-1, keepdim=True)
    tf2 /= tf2.norm(dim=-1, keepdim=True)
    text_sim = torch.cosine_similarity(tf1, tf2).item()




if text_sim > 0.9:
    image_sim = torch.cosine_similarity(emb1, emb2).item()
    print(f"Сходство изображений: {image_sim * 100:.2f}%")
    print(f"Сходство описаний: {text_sim * 100:.2f}%")
else:
    print("Описание слишком разное. Картинки не сравниваются.")
    print(f"Сходство описаний: {text_sim * 100:.2f}%")
