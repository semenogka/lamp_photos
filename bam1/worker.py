
import torch
import clip
from PIL import Image
import requests
from io import BytesIO
from PyQt6.QtCore import QThread, pyqtSignal
import base64
from .helper import is_number
import traceback
device = "cuda" if torch.cuda.is_available() else "cpu"
model, preprocess = clip.load("ViT-B/32", device)


labels = [
    "lamp in the shape of a circle",
    "lamp with a round-shaped body",
    "flush‑mount disc lamp",

    "rectangular lamp",
    "lamp with an elongated rectangular shape",
    "thin rectangular ceiling lamp",

    "square lamp",
    "lamp with a square body",
    "compact square lamp",

    "lamp in the shape of a sphere",
    "spherical pendant lamp",
    "voluminous lamp resembling a sphere",

    "lamp in the shape of a cube",
    "cube-shaped lamp body",
    "geometrically strict cube-shaped lamp",

    "long elongated lamp",
    "narrow lamp in the shape of a strip",
    "linear elongated lamp",

    "lamp with an angular shape",
    "lamp with sharp angles and edges",
    "lamp with a sharply geometric body",

    "conical lampshade",
    "lamp with a bell-shaped shade",
    "lampshade tapering upwards",

    "cylindrical lamp",
    "vertical cylinder-shaped lamp",
    "lamp body in the shape of a tube"
]

text_tokens = clip.tokenize(labels).to(device)
with torch.no_grad():
    text_features = model.encode_text(text_tokens)
    text_features /= text_features.norm(dim=-1, keepdim=True)

def get_image_features_and_label(image: Image.Image):
    embs, best_text_features = [], []
    mirrored = image.transpose(Image.FLIP_LEFT_RIGHT)

    
    for r in [0, 90, 180, 270]:
        with torch.no_grad():
            orig = image.rotate(r, expand=True)
            mir = mirrored.rotate(r, expand=True)

            imageOrig = preprocess(orig).unsqueeze(0).to(device)
            imageMir = preprocess(mir).unsqueeze(0).to(device)
            for img_tensor in [imageOrig, imageMir]:
                image_features = model.encode_image(img_tensor)
                image_features /= image_features.norm(dim=-1, keepdim=True)
                similarities = (image_features @ text_features.T).squeeze(0)
                best_idx = similarities.argmax().item()
                best_text_features.append(text_features[best_idx].unsqueeze(0))
                embs.append(image_features)

    return embs, best_text_features




class Worker(QThread):
    finished = pyqtSignal(list)
    
    def __init__(self, img, cats, widgth, height, abajurcolors, metalcolors, abajurmaterial, metalmaterial, _cache, abajur_colors_dict, metal_colors_dict, abajur_materials_dict, armatur_materials_dict):
        super().__init__()
        self.img = img.convert("RGB")  
        self.cats = cats
        self.widgth = widgth
        self.height = height
        self.abajurcolors = abajurcolors
        self.metalcolors = metalcolors
        self.abajurmaterial = abajurmaterial
        self.metalmaterial = metalmaterial
        self.cache = _cache
        self.abajur_colors_dict = abajur_colors_dict
        self.metal_colors_dict = metal_colors_dict
        self.abajur_materials_dict = abajur_materials_dict
        self.armatur_materials_dict = armatur_materials_dict
    def run(self):
        try:
            result = []

            image = self.img.convert("RGB")
            embs, labels = get_image_features_and_label(image)
        
            for category in self.cats:
                if not isinstance(self.cache, dict):
                    print("ОШИБКА: _cache не словарь, а", type(self.cache))
                    break

                if category not in self.cache:
                    print(f"ОШИБКА: Нет ключа '{category}' в _cache, available:", self.cache.keys())
                    continue

                cache_entry = self.cache[category]
                if not isinstance(cache_entry, dict):
                    print(f"ОШИБКА: Для категории '{category}' в _cache лежит {type(cache_entry)}, а не dict")
                    continue

                # Убедимся, что есть всё, что нужно:
                for needed in ("data","text_embs","emb_array"):
                    if needed not in cache_entry:
                        print(f"ОШИБКА: В cache_entry[{category!r}] нет поля '{needed}'")
                        break
                else:
                    # Только если все ключи на месте, продолжаем
                    data      = cache_entry["data"]
                    text_embs = cache_entry["text_embs"]
                    emb_array = cache_entry["emb_array"]

                for i in range(len(emb_array)):
                        minres = []  
                        # for embs_img, labels_img in zip(embs, labels):
                        #     for tf1, emb_user in zip(labels_img, embs_img):
                        for embs_img, labels_img in zip(embs, labels):
                            for tf1, emb_user in zip(labels_img, embs_img):
                                text_sim = torch.cosine_similarity(tf1, text_embs[i], dim=-1)
                                if text_sim > 0.9:
                                    item = data[i]
                                    
                                    
                                    if self.widgth.is_enabled():
                                        mind, maxd = self.widgth.get_values()
                                        widgth_ok = (
                                            (is_number(item.get("widght")) and mind <= int(item["widght"]) <= maxd) or
                                            (is_number(item.get("diameter")) and mind <= int(item["diameter"]) <= maxd) or
                                            (is_number(item.get("length")) and mind <= int(item["length"]) <= maxd)
                                        )
                                        if not widgth_ok:
                                            continue

                                    if self.height.is_enabled():
                                        mind, maxd = self.height.get_values()
                                        if not (is_number(item.get("height")) and mind <= int(item["height"]) <= maxd):
                                            continue
                                    if self.abajurcolors.is_enabled():
                                        colorCheck = False
                                        if item.get("abajurcolor") and item["abajurcolor"] != "None":
                                            if item["abajurcolor"] == "None":
                                                print(item["abajurcolor"])
                                            colors = self.abajurcolors.selected_colors()
                                            
                                            for key in colors:
                                                ncolors = self.abajur_colors_dict[key]
                                                for col in ncolors:
                                                    if col == item["abajurcolor"]:
                                                        colorCheck = True
                                        if colorCheck == False:
                                                continue

                                    if self.metalcolors.is_enabled():
                                        colorCheck = False
                                        
                                        if item.get("metalcolor") and item["metalcolor"] != "None":
                                            colors = self.metalcolors.selected_colors()
                                            
                                            for key in colors:
                                                ncolors = self.metal_colors_dict[key]
                                                for col in ncolors:
                                                    if col == item["metalcolor"]:
                                                        print(item["name"])
                                                        colorCheck = True
                                        if colorCheck == False:
                                            continue

                                    if self.abajurmaterial.is_enabled():
                                        materialCheck = False
                                        if item.get("abajurmaterial") and item["abajurmaterial"] != "None":
                                            materials = self.abajurmaterial.selected_colors()
                                            
                                            for key in materials:
                                                nmaterials = self.abajur_materials_dict[key]
                                                for col in nmaterials:
                                                    if col == item["abajurmaterial"]:
                                                        print(item["name"])
                                                        materialCheck = True
                                        if materialCheck == False:
                                            continue
                                    if self.metalmaterial.is_enabled():
                                        materialCheck = False
                                        if item.get("armaturmaterial") and item["armaturmaterial"] != "None":
                                            materials = self.metalmaterial.selected_colors()
                                            
                                            for key in materials:
                                                nmaterials = self.armatur_materials_dict[key]
                                                for col in nmaterials:
                                                    if col == item["armaturmaterial"]:
                                                        print(item["name"])
                                                        materialCheck = True
                                        if materialCheck == False:
                                            continue
                                    # Если фильтры прошли, считаем схожесть по изображению
                                    sims = []
                                    for emb in emb_array[i]:
                                        image_sim = torch.cosine_similarity(emb_user, emb, dim=-1)
                                        sims.append(image_sim)
                                    minres.append({
                                        "sim": float(max(sims) * 100),
                                        "link": item['link'],
                                        "name": item['name']
                                    })
                        if minres:
                            best = max(minres, key=lambda x: x['sim'])
                            
                            result.append(best)
                
                
            
            
            
            result_sorted = sorted(result, key=lambda x: x["sim"], reverse=True)

            res_results = []
            for r in result_sorted[:150]:
                url = f"https://raw.githubusercontent.com/semenogka/lamp_photos/main/allimgs/{r['name']}"
                response = requests.get(url)
                if response.status_code == 200:
                    img_bytes = BytesIO(response.content)
                    base64_data = base64.b64encode(img_bytes.getvalue()).decode('utf-8')
                    
                    res_results.append({
                        "base64": base64_data,
                        "link": r['link'],
                        "sim": r['sim']
                    })
                

            self.finished.emit(res_results)
            
        except Exception as e:
            print("‼️ Ошибка в run():", e)
            traceback.print_exc()