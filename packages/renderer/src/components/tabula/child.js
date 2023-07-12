export default childWindow = {
    main: {
        window: window,
        document: window.document,
        alert: (message) => {
            window.alert(message)
        },
    },
}
