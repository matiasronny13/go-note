const app = (function () {    
    const bindingDialogButtons = () => {
        const dialog = document.querySelector("dialog");
        const showButton = document.getElementById("deleteButton")
        const cancelButton = document.getElementById("cancelButton")
        
        if (!!showButton) {
            showButton.addEventListener("click", () => {
                dialog.showModal();
            });
        }

        cancelButton.addEventListener("click", () => {
            dialog.close();
        });
    }

    htmx.onLoad(function(elt) {
        if (htmx.findAll(elt, "#contentRoot").length > 0) {
            bindingDialogButtons();
        }
    })

    const markedOptions = {
        pedantic: false,
        gfm: true,
        breaks: true
    };

    return {
        closeDialog: () => {dialog.close();},
        markdownMode: (content) => {
            if (!!content) {
                placeholder = document.createElement("span");
                placeholder.innerHTML = content;            
                document.getElementById("xData")._x_dataStack[0].isEdit = false;
                return marked.parse([].slice.call(placeholder.childNodes).map(a => a.textContent).join('\n'), markedOptions);
            } else {
                return "";
            }
        },
        editMode: (element, content) => {
            element.innerHTML = content
            document.getElementById("xData")._x_dataStack[0].isEdit = true;
        }
    }
})();