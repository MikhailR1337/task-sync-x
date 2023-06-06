(function() {
    const forms = document.querySelectorAll('form');
    for (const form of forms) {
        const childs = form.childNodes;
        let input = null;
        for (const child of childs) {
            if (child.type === 'hidden') {
                input = child;
                break;
            }
        }
        if (!input) {
            return;
        }
    
        const reload = async () => {
            location.reload()
        }
        const sendData = async () => {
            const formData = new FormData(form);
            try {
                await fetch(form.action, {
                    method: input.value,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(Object.fromEntries(formData)),
                })
                await reload();
            } catch(e) {
                console.error(e)
            }
    
        } 
        form.addEventListener('submit', (event) => {
            event.preventDefault();
            sendData()
        });

    }
})()