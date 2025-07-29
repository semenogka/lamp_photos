import sys
from PyQt6.QtCore import Qt
from PyQt6.QtGui import QImage, QPixmap
from PyQt6.QtWidgets import QApplication, QLabel, QMainWindow, QVBoxLayout, QWidget, QPushButton, QMessageBox, QHBoxLayout


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()

        self.setWindowTitle("Вставка до 3 скриншотов из буфера обмена")
        self.setGeometry(100, 100, 900, 400)

        # Список QLabel для отображения изображений
        self.image_labels = [QLabel(self) for _ in range(3)]
        for label in self.image_labels:
            label.setAlignment(Qt.AlignmentFlag.AlignCenter)
            label.setFixedSize(280, 280)
            label.setStyleSheet("border: 1px solid black;")

        # Кнопка вставки
        self.paste_button = QPushButton("Вставить скриншот", self)
        self.paste_button.clicked.connect(self.paste_image)

        # Макеты
        images_layout = QHBoxLayout()
        for label in self.image_labels:
            images_layout.addWidget(label)

        main_layout = QVBoxLayout()
        main_layout.addLayout(images_layout)
        main_layout.addWidget(self.paste_button)

        container = QWidget()
        container.setLayout(main_layout)
        self.setCentralWidget(container)

        # Счетчик вставленных изображений
        self.inserted_images = 0

    def paste_image(self):
        if self.inserted_images >= 3:
            QMessageBox.information(self, "Информация", "Максимум 3 изображения уже вставлено.")
            return

        clipboard = QApplication.clipboard()
        mime_data = clipboard.mimeData()

        if mime_data.hasImage():
            image = QImage(mime_data.imageData())
            if not image.isNull():
                pixmap = QPixmap(image).scaled(
                    self.image_labels[self.inserted_images].width(),
                    self.image_labels[self.inserted_images].height(),
                    Qt.AspectRatioMode.KeepAspectRatio,
                    Qt.TransformationMode.SmoothTransformation
                )
                self.image_labels[self.inserted_images].setPixmap(pixmap)
                self.inserted_images += 1
            else:
                QMessageBox.warning(self, "Ошибка", "Не удалось получить изображение из буфера обмена.")
        else:
            QMessageBox.warning(self, "Ошибка", "В буфере обмена нет изображения.")


if __name__ == "__main__":
    app = QApplication(sys.argv)
    window = MainWindow()
    window.show()
    sys.exit(app.exec())
