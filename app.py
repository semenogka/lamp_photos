import os
import torch
import clip
from PIL import Image, ImageQt
import json
import requests
from io import BytesIO
from PyQt6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QPushButton, QListWidget, QLabel,
    QAbstractItemView, QListWidgetItem, QTextBrowser, QMessageBox
)
from PyQt6.QtCore import Qt
from PyQt6.QtGui import QPixmap, QGuiApplication, QKeySequence, QAction
import sys
import base64


def resource_path(relative_path):
    if hasattr(sys, '_MEIPASS'):
        return os.path.join(sys._MEIPASS, relative_path)
    return os.path.join(os.path.abspath("."), relative_path)

# Устройство и загрузка модели
device = "cuda" if torch.cuda.is_available() else "cpu"
model, preprocess = clip.load("ViT-B/32", device)



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


base_path = resource_path(os.path.join("alldata", "ndata"))

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



class MainWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Поиск светильников")
        self.resize(1000, 1000)

        self.layout = QVBoxLayout()
        self.setLayout(self.layout)
        
        paste_action = QAction(self)
        paste_action.setShortcut(QKeySequence.StandardKey.Paste)
        paste_action.triggered.connect(self.paste_image)
        self.addAction(paste_action)
        self.setFocusPolicy(Qt.FocusPolicy.StrongFocus)

        self.list_categories = QListWidget()
        self.list_categories.setSelectionMode(QAbstractItemView.SelectionMode.MultiSelection)
        self.layout.addWidget(QLabel("Выберите категории:"))
        self.layout.addWidget(self.list_categories)

        for cat in buttons_files.keys():
            item = QListWidgetItem(cat)
            self.list_categories.addItem(item)

        self.image_preview = QLabel()
        self.layout.addWidget(self.image_preview)

        self.btn_search = QPushButton("Запустить поиск")
        self.btn_search.clicked.connect(self.run_search)     
        self.layout.addWidget(self.btn_search)      

        self.results_browser = QTextBrowser()
        self.layout.addWidget(self.results_browser)


        self.results_browser.setOpenExternalLinks(True)
        self.img = 0

    def paste_image(self):
        clipboard = QGuiApplication.clipboard()
        mime_data = clipboard.mimeData()

        if mime_data.hasImage():
            qimage = clipboard.image()  # QImage
            pixmap = QPixmap.fromImage(qimage)
            self.image_preview.setPixmap(pixmap.scaled(300, 300, Qt.AspectRatioMode.KeepAspectRatio))

            # Преобразуем QImage → PIL.Image правильно
            image_pil = ImageQt.fromqimage(qimage).convert("RGB")
            self.img = image_pil
        else:
            QMessageBox.warning(self, "Ошибка", "Буфер обмена не содержит изображения")
        print("вставил")

            

    def run_search(self):
        print("начал поиск")
        self.results_browser.clear()
        if self.img == 0:
            QMessageBox.warning(self, "Ошибка", "Сначала загрузите изображения")
            return
        cats = [item.text() for item in self.list_categories.selectedItems()]
        if not cats:
            QMessageBox.warning(self, "Ошибка", "Выберите хотя бы одну категорию")
            return
        
        self.results_browser.clear()
        self.results_browser.append("Вычисляем эмбеддинги изображений...")

        result = []
        # all_embs, all_labels = [], []
        
        image = self.img.convert("RGB")
        embs, labels = get_image_features_and_label(image)
        print("собрал эмб")

        for category in cats:
            print("начало категории")
            QApplication.processEvents()
            clicked_file = buttons_files[category]
            with open(clicked_file, "r", encoding="utf-8") as f:
                data = json.load(f)
            for item in data:
                QApplication.processEvents()
                print("начало файла")
                minRes = []
                text2 = clip.tokenize([item["label"]]).to(device)
                with torch.no_grad():
                    tf2 = model.encode_text(text2)
                    tf2 /= tf2.norm(dim=-1, keepdim=True)
                emb = torch.tensor(item['emb'], dtype=torch.float32).to(device)
                for embs_img, labels_img in zip(embs, labels):
                    QApplication.processEvents()
                    for tf1, emb_user in zip(labels_img, embs_img):
                        print("сравнение")
                        text_sim = torch.cosine_similarity(tf1, tf2).item()
                        if text_sim > 0.95:
                            image_sim = torch.cosine_similarity(emb_user, emb, dim=-1).item()
                            minRes.append(image_sim)
                if minRes:
                    print("минрес")
                    good_sim = max(minRes)
                    if good_sim > 0.7:
                        result.append({
                            "sim": good_sim,
                            "link": item['link'],
                            "name": item['name']
                        })

        result_sorted = sorted(result, key=lambda x: x["sim"], reverse=True) 
        self.results_browser.append("\nЛучшие совпадения:")
        print("сортед")
        for r in result_sorted[:60]:
            QApplication.processEvents()
            url = f"https://raw.githubusercontent.com/semenogka/lamp_photos/main/allimgs/{r['name']}"
            response = requests.get(url)
            if response.status_code == 200:
                img_bytes = BytesIO(response.content)
                base64_data = base64.b64encode(img_bytes.getvalue()).decode('utf-8')
            html = f'''
            <div style="margin-bottom:15px;">
                <img src="data:image/jpeg;base64,{base64_data}" width="100" style="vertical-align: middle; margin-right: 10px;">
                <a href="{r["link"]}">{r["link"]}</a><br>
                <span>Похожесть: {r["sim"]:.3f}</span>
            </div>
        '''
            self.results_browser.append(html)


app = QApplication(sys.argv)
win = MainWindow()
win.show()
sys.exit(app.exec())

