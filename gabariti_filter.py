from PyQt6.QtWidgets import (
     QWidget,  QLabel, QCheckBox, QDoubleSpinBox, QHBoxLayout
)

class GabaritiFilter(QWidget):
    def __init__(self, name: str):
        super().__init__()
        self.name = name
        self.init_ui()

    def init_ui(self):
        layout = QHBoxLayout()
        self.checkbox = QCheckBox(self.name)
        self.checkbox.stateChanged.connect(self.toggle_inputs)

        self.min_input = QDoubleSpinBox()
        self.min_input.setRange(0, 10000)
        self.min_input.setSuffix(" мм")
        self.min_input.setDecimals(2)
        self.min_input.setEnabled(False)

        self.max_input = QDoubleSpinBox()
        self.max_input.setRange(0, 10000)
        self.max_input.setSuffix(" мм")
        self.max_input.setDecimals(2)
        self.max_input.setEnabled(False)

        layout.addWidget(self.checkbox)
        layout.addWidget(QLabel("от:"))
        layout.addWidget(self.min_input)
        layout.addWidget(QLabel("до:"))
        layout.addWidget(self.max_input)

        self.setLayout(layout)

    def toggle_inputs(self, state):
        enabled = state == 2
        self.min_input.setEnabled(enabled)
        self.max_input.setEnabled(enabled)

    def is_enabled(self):
        return self.checkbox.isChecked()

    def get_values(self):
        return (self.min_input.value(), self.max_input.value())