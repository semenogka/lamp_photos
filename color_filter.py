
from PyQt6.QtWidgets import (
    QWidget, QVBoxLayout, QListWidget, QListWidgetItem, QCheckBox
)
from PyQt6.QtCore import Qt

class ColorFilter(QWidget):
    def __init__(self, name: str, color_list: dict[str, list[str]]):
        super().__init__()
        self.name = name
        self.color_list = color_list
        self.init_ui()

    def init_ui(self):
        layout = QVBoxLayout()
        self.checkbox = QCheckBox(self.name)
        self.checkbox.stateChanged.connect(self.toggle_input)

        self.list_widget = QListWidget()
        self.list_widget.setEnabled(False)
        for color in sorted(self.color_list.keys()):
            item = QListWidgetItem(color)
            item.setFlags(item.flags() | Qt.ItemFlag.ItemIsUserCheckable)
            item.setCheckState(Qt.CheckState.Unchecked)
            self.list_widget.addItem(item)
        
        layout.addWidget(self.checkbox)
        layout.addWidget(self.list_widget)
        self.setLayout(layout)

    def toggle_input(self, state):
        self.list_widget.setEnabled(state)

    def selected_colors(self) -> list[str]:
        """Возвращает список текстов отмеченных элементов, 
        если фильтр включён, иначе пустой список."""
        if not self.is_enabled():
            return []
        result = []
        for i in range(self.list_widget.count()):
            item = self.list_widget.item(i)
            if item.checkState() == Qt.CheckState.Checked:
                result.append(item.text())
        return result
    
    def is_enabled(self):
        return self.checkbox.isChecked()