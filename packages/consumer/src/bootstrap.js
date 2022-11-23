import {
    BIcon,
    BIconDashSquareFill,
    BIconPlusSquareFill,
    BInputGroup,
    ButtonPlugin,
    FormPlugin,
    FormSelectPlugin,
    ModalPlugin,
    LayoutPlugin,
    OverlayPlugin
} from "bootstrap-vue";

export default {
    install: (app, options) => {
        app.component("BIcon", BIcon)
        app.component("BIconDashSquareFill", BIconDashSquareFill)
        app.component("BIconPlusSquareFill", BIconPlusSquareFill)
        app.component("BInputGroup", BInputGroup)
        app.use(ButtonPlugin)
        app.use(FormPlugin)
        app.use(FormSelectPlugin)
        app.use(LayoutPlugin)
        app.use(ModalPlugin)
        app.use(OverlayPlugin)
    }
}
