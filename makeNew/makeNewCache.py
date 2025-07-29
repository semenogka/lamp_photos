import json
import os
import torch
import clip
import sys

def resource_path(relative_path):
    if hasattr(sys, '_MEIPASS'):
        return os.path.join(sys._MEIPASS, relative_path)
    return os.path.join(os.path.abspath("."), relative_path)

base_path = resource_path(os.path.join("alldata", "ndata"))

buttons_files = {
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

files = [
    "bra",
    "fasad",
    "groont",
    "land",
    "lystri",
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
device = "cuda" if torch.cuda.is_available() else "cpu"
model, preprocess = clip.load("ViT-B/32", device)

def get_cache():
    for id, category in enumerate(buttons_files, start=0):
        clicked_file = buttons_files[category]
        with open(clicked_file, "r", encoding="utf-8") as f:
            data = json.load(f)

        text_labels = [item["label"] for item in data]
        text_tokens = clip.tokenize(text_labels).to(device)

        text_embs_list = []
        BATCH_SIZE = 512
        with torch.no_grad():
            for i in range(0, len(text_tokens), BATCH_SIZE):
                batch = text_tokens[i : i + BATCH_SIZE]
                batch_embs = model.encode_text(batch)
                batch_embs /= batch_embs.norm(dim=-1, keepdim=True)
                text_embs_list.append(batch_embs)
        text_embs = torch.cat(text_embs_list, dim=0)  

        emb_list = [torch.tensor(item["emb"], dtype=torch.float32) for item in data]
        emb_array = torch.stack(emb_list).to(device)  

        entry = {
            "data":     data,
            "text_embs": text_embs,
            "emb_array": emb_array
        }
        print("ready", category)
        torch.save(entry, f"cache/{files[id]}.pt")


get_cache()