import { store } from '../../store'
import { loadManuscripts } from './manuscript'

const originalColorOptions = [
    {
        value: 0,
        text: 'Svart text på vit bakgrund',
        colors: {
            foreground: '#ffffff',
            background: '#000000',
        },
    },
    {
        value: 1,
        text: 'Vit text på svart bakgrund',
        colors: {
            foreground: '#000000',
            background: '#ffffff',
        },
    },
    {
        value: 2,
        text: 'Gul text på svart bakgrund',
        colors: {
            foreground: '#ffff00',
            background: '#000000',
        },
    },
    {
        value: 3,
        text: 'Svart text på gul bakgrund',
        colors: {
            foreground: '#000000',
            background: '#ffff00',
        },
    },
    {
        value: 4,
        text: 'Gul text på blå bakgrund',
        colors: {
            foreground: '#ffff00',
            background: '#0000ff',
        },
    },
    {
        value: 5,
        text: 'Blå text på gul bakgrund',
        colors: {
            foreground: '#0000ff',
            background: '#ffff00',
        },
    },
]
var colorOptions = []
function loadTextSettings() {
    let settings = store.state.settings
    let font = settings.font
    let behaviour = settings.behaviour || {}
    if (behaviour == null) {
        if (font.capitalizeOnNewLine != undefined) {
            behaviour.capitalizeOnNewLine = font.capitalizeOnNewLine
        } else {
            behaviour.capitalizeOnNewLine = true
        }
    }
    colorOptions = originalColorOptions.slice()
    if (font.customColors.valid && colorOptions) {
        colorOptions.push({
            value: 6,
            text: 'Eget färgschema',
            colors: {
                foreground: font.customColors.foreground,
                background: font.customColors.background,
            },
        })
    }

    return {
        font: {
            fontFamily: font.family,
            fontSize: font.size + 'px',
            backgroundColor: font.background,
            color: font.foreground,
            lineHeight: font.lineHeight,
            colorID: font.colorID,
            customColors: font.customColors,
            marginTop: 0 + 'px',
            marginRight: 0 + 'px',
            marginBottom: 0 + 'px',
            marginLeft: 0 + 'px',
        },
        behaviour: {
            capitalizeOnNewLine:
                behaviour.capitalizeOnNewLine || font.capitalizeOnNewLine,
        },
    }
}
function saveTextSettings(current) {
    let settings = {
        font: {
            family: current.font.fontFamily,
            size: parseInt(current.font.fontSize),
            lineHeight: current.font.lineHeight,
            background: current.font.backgroundColor,
            foreground: current.font.color,
            colorID: current.font.colorID,
            customColors: current.font.customColors,
            margins: {
                top: parseInt(current.font.marginTop),
                right: parseInt(current.font.marginRight),
                bottom: parseInt(current.font.marginBottom),
                left: parseInt(current.font.marginLeft),
            },
        },
        behaviour: {
            capitalizeOnNewLine: current.behaviour.capitalizeOnNewLine,
        },
    }
    store.commit('setFontSettings', settings)
}

function changeTextSize(inc, size) {
    size = parseInt(size)
    size = size + (inc ? 4 : -4)
    if (size <= 20) {
        size = 20
    } else if (size >= 200) {
        size = 200
    }
    return size + 'px'
}

function changeColor(colorID) {
    if (colorID != null) {
        colorID++
        if (colorID > colorOptions.length - 1) {
            colorID = 0
        }
    } else {
        colorID = 0
    }
    return {
        background: colorOptions[colorID].colors.background,
        foreground: colorOptions[colorID].colors.foreground,
        colorID: colorID,
    }
}

function initText() {
    loadManuscripts(store.getters.selectedManuscripts)
}

export default {
    initText,
    loadTextSettings,
    saveTextSettings,
    changeTextSize,
    changeColor,
}
