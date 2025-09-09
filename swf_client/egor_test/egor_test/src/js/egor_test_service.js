import appCtx from 'js/appCtxService';

export function updatePMI(editor,elementName, attributes) {
  const selection = editor.model.document.selection;
  const selectedElement = selection.getSelectedElement();

  editor.model.change(writer => {
    const newElement = writer.createElement(elementName, attributes);
    writer.insert(newElement, writer.createPositionAfter(selectedElement));
    writer.remove(selectedElement);
  });
}

export async function insertSelectedPMI(editor,elementName, selectedPMI){
  editor.model.change(writer => {
    let upper,lower = ""
    if (selectedPMI.props.upperDelta) upper = selectedPMI.props.upperDelta.value
    if (selectedPMI.props.lowerDelta) lower = selectedPMI.props.lowerDelta.value
    if (upper.length && (upper[0] === "+" || upper[0] === "-")) upper = upper = upper.slice(1);
    if (lower.length && (lower[0] === "+" || upper[0] === "-")) lower = lower = lower.slice(1);
    const newElement = writer.createElement(elementName, {
      diameter: selectedPMI.props.value.value,
      upper: upper,
      lower: lower,
      pmiName: selectedPMI.props.name.value
    });
    editor.model.insertContent(newElement);
  });
}

