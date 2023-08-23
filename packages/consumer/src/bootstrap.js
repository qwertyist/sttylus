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
  ListGroupPlugin,
  ModalPlugin,
  LayoutPlugin,
  OverlayPlugin,
  TabsPlugin,
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
    app.use(ListGroupPlugin)
    app.use(ModalPlugin)
    app.use(OverlayPlugin)
    app.use(TabsPlugin)
  },
}
