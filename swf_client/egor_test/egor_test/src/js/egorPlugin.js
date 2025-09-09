import { ckeditor5ServiceInstance } from 'js/Arm0CkeditorServiceInstance';
//import { ckeditor5ServiceInstance } from 'js/wiEditor.service';
import commandPanelService from 'js/commandPanel.service';

export default class EgorPlugin extends ckeditor5ServiceInstance.Plugin {
    init() {
        const editor = this.editor;

        editor.model.schema.register('egorTag', {
            allowWhere: '$text',
            isObject: true,
            isInline: true,
            allowAttributes: ['pmiName','diameter', 'upper', 'lower']
        });


        editor.ui.componentFactory.add('egorPlugin', locale => {
            const view = new ckeditor5ServiceInstance.ButtonView(locale);
            view.set({
                label: 'Insert Egor',
                withText: true,
                tooltip: true
            });

            // add template object
            // view.on('execute', () => {
            //     editor.model.change(writer => {
            //         const diameterElement = writer.createElement('egorTag', {
            //             diameter: '40',
            //             upper: '0.2',
            //             lower: '0.1'
            //         });
            //         editor.model.insertContent(diameterElement);
            //     });
            // });

            view.on('execute', () => {
                editor.model.change(writer => {
                    commandPanelService.activateCommandPanel('EgorPmiPanelMain', 'aw_toolsAndInfo', {
                    editor: editor,
                    elementName: 'egorTag',
                });
                });
            });

            return view;

        })


        // Downcast: модель → HTML
        editor.conversion.for('downcast').elementToElement({
            model: 'egorTag',
            view: (modelItem, { writer }) => {
                const diameter = modelItem.getAttribute('diameter');
                const upper = modelItem.getAttribute('upper');
                const lower = modelItem.getAttribute('lower');
                const pmiName = modelItem.getAttribute('pmiName');
                

                const container = writer.createContainerElement('span', {
                    class: 'egor-tag',
                    'data-diameter': diameter,
                    'data-upper': upper,
                    'data-lower': lower,
                    'data-pmi-name':pmiName,
                    style: 'display:inline-block; padding:2px 4px; font-family:monospace;'
                });

                const diameterSpan = writer.createContainerElement('span', {
                    style: 'font-weight:bold;'
                });
                writer.insert(writer.createPositionAt(diameterSpan, 0), writer.createText(`⌀${diameter}`));
                let upperSpan,lowerSpan = null
                if (upper && upper.length){
                    upperSpan = writer.createContainerElement('span', {
                        style: 'font-size:0.8em; vertical-align:super;'
                    });
                    writer.insert(writer.createPositionAt(upperSpan, 0), writer.createText(`+${upper}`));
                }
                if (lower && lower.length){
                    lowerSpan = writer.createContainerElement('span', {
                        style: 'font-size:0.8em; vertical-align:sub;'
                    });
                    writer.insert(writer.createPositionAt(lowerSpan, 0), writer.createText(`-${lower}`));
                }
                writer.insert(writer.createPositionAt(container, 0), diameterSpan);
                if (upperSpan) writer.insert(writer.createPositionAt(container, 'end'), upperSpan);
                if (lowerSpan) writer.insert(writer.createPositionAt(container, 'end'), lowerSpan);

                return ckeditor5ServiceInstance.toWidget(container, writer);
            }
        });

        // Upcast: HTML → модель
        editor.conversion.for('upcast').elementToElement({
            view: {
                name: 'span',
                class: 'egor-tag',
                attributes: {
                    'data-diameter': true,
                    'data-upper': true,
                    'data-lower': true,
                    'data-pmi-name':true
                }
            },
            model: (viewElement, { writer }) => {
                return writer.createElement('egorTag', {
                    diameter: viewElement.getAttribute('data-diameter'),
                    upper: viewElement.getAttribute('data-upper'),
                    lower: viewElement.getAttribute('data-lower'),
                    pmiName: viewElement.getAttribute('data-pmi-name'), 
                });
            }
        });


        // Обработка клика по элементу
        editor.editing.view.document.on('click', (evt, data) => {
            const viewElement = data.target;

            const container = viewElement.findAncestor(el => {
                return el.is('element') && el.hasClass('egor-tag');
            });

            if (container) {
                const diameter = container.getAttribute('data-diameter')
                const upper = container.getAttribute('data-upper')
                const lower = container.getAttribute('data-lower')
                const pmiName = container.getAttribute('data-pmi-name')
                commandPanelService.activateCommandPanel('egorEditPanel', 'aw_toolsAndInfo', {
                    editor: editor,
                    elementName: 'egorTag',
                    attributes: {
                        diameter,
                        upper,
                        lower,
                        pmiName
                    }
                });
            }
        })
    }
}
