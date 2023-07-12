import {
  BadgePlugin,
  BIcon,
  BIconDashSquareFill,
  BIconPlusSquareFill,
  BInputGroup,
  ButtonPlugin,
  FormInputPlugin,
  FormPlugin,
  FormSelectPlugin,
  ModalPlugin,
  LayoutPlugin,
  OverlayPlugin,
} from 'bootstrap-vue'

export default {
  install: (app, options) => {
    app.component('BIcon', BIcon)
    app.component('BIconDashSquareFill', BIconDashSquareFill)
    app.component('BIconPlusSquareFill', BIconPlusSquareFill)
    app.component('BInputGroup', BInputGroup)
    app.use(BadgePlugin)
    app.use(ButtonPlugin)
    app.use(FormPlugin)
    app.use(FormInputPlugin)
    app.use(FormSelectPlugin)
    app.use(LayoutPlugin)
    app.use(ModalPlugin)
    app.use(OverlayPlugin)
  },
}
