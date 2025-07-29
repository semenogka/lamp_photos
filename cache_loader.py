from helper import resource_path
from helper import file_sha256
import hashlib
import torch
import requests
import os

_cache = {}

def get_files(buttons_files):
    os.makedirs(resource_path("cache"), exist_ok=True) 

    for file in buttons_files:
        filename = os.path.basename(buttons_files[file])
        raw_url = f"https://raw.githubusercontent.com/semenogka/lamp_photos/main/cache/{filename}"

        print(f"Загрузка: {file}")
        res = requests.get(raw_url)
        if res.status_code == 200:
            local_path = resource_path(os.path.join("cache", filename))
            remote_hash = hashlib.sha256(res.content).hexdigest()
            local_hash = file_sha256(local_path)
            if local_hash != remote_hash or not os.path.exists(local_path):
              with open(local_path, "wb") as f:
                f.write(res.content)
        else:
            print(f"Ошибка загрузки {file}: HTTP {res.status_code}")

def load_cache(buttons_files):
    global _cache
    for category in buttons_files:
      pt_path = resource_path(buttons_files[category])
      print(pt_path)
      if not os.path.exists(pt_path):
            print(f"[!] Не найден .pt-файл для '{category}': {pt_path}")
            continue

      try:
            entry = torch.load(pt_path, map_location="cpu")
            _cache[category] = entry
      except Exception as e:
            print(f"[!] Ошибка torch.load({pt_path}): {e}")
            continue

        # Проверяем, что структуру выложили правильно
      if not all(k in entry for k in ("data", "text_embs", "emb_array")):
            print(f"[!] Неверная структура в {pt_path}, ключи: {entry.keys()}")
            continue
    return _cache