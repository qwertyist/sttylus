import {
    BadgePlugin,
    BIcon,
    BIconBrush,
    BIconDownload,
    BIconExclamationCircleFill,
    BIconTrash,
    BInputGroup,
    BPopover,
    BSpinner,
    VBTooltip,
    ButtonPlugin,
    BInputGroupAppend,
    CardPlugin,
    FormPlugin,
    FormGroupPlugin,
    FormInputPlugin,
    InputGroupPlugin,
    JumbotronPlugin,
    LayoutPlugin,
    ListGroupPlugin,
    ModalPlugin,
    NavbarPlugin,
    OverlayPlugin,
    PaginationPlugin,
    SidebarPlugin,
    TablePlugin,
    TabsPlugin,
    FormRadioPlugin,
    FormCheckboxPlugin,
    FormSelectPlugin,
    FormFilePlugin,
    FormTextareaPlugin,
} from 'bootstrap-vue'

export default {
    install: (app, options) => {
        app.use(BadgePlugin)
        app.use(ButtonPlugin)
        app.use(CardPlugin)
        app.use(FormCheckboxPlugin)
        app.use(FormFilePlugin)
        app.use(FormTextareaPlugin)
        app.use(FormGroupPlugin)
        app.use(FormInputPlugin)
        app.use(FormPlugin)
        app.use(FormRadioPlugin)
        app.use(FormSelectPlugin)
        app.use(InputGroupPlugin)
        app.use(JumbotronPlugin)
        app.use(LayoutPlugin)
        app.use(ListGroupPlugin)
        app.use(ModalPlugin)
        app.use(NavbarPlugin)
        app.use(OverlayPlugin)
        app.use(PaginationPlugin)
        app.use(SidebarPlugin)
        app.use(TablePlugin)
        app.use(TabsPlugin)
        app.component('BIcon', BIcon)
        app.component('BIconBrush', BIconBrush)
        app.component('BIconDownload', BIconDownload)
        app.component('BIconExclamationCircleFill', BIconExclamationCircleFill)
        app.component('BIconTrash', BIconTrash)
        app.component('BInputGroup', BInputGroup)
        app.component('BInputGroupAppend', BInputGroupAppend)
        app.component('BPopover', BPopover)
        app.component('BSpinner', BSpinner)
        app.directive('b-tooltip', VBTooltip)
    },
}
