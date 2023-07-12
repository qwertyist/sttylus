<template>
    <div>
        <ColorPicker ref="colorPicker" />
        <b-row>
            <b-col cols="4">
                <div id="fontSettings">
                    <b-form @change="onChange" @reset="onReset">
                        <!--    <b-form-group label="Typsnitt">
        <b-form-select name="fontFamily" id="fontFamily" v-model="form.family.selected">
          <option
            v-for="(fontFamily, indexOpt) in form.family.options"
            :key="indexOpt"
            :value="fontFamily"
          >{{fontFamily}}</option>
        </b-form-select>
      </b-form-group>
            -->
                        <!--    <b-form-group label="Textbeteende">
              <b-form-checkbox v-model="form.capitalizeOnNewLine">
                Stor bokstav vid ny rad
              </b-form-checkbox>
            </b-form-group>
            -->
                        <b-row>
                            <b-col cols="3">
                                <div class="pt-2">Typsnitt</div>
                            </b-col>
                            <b-col>
                                <b-form-select
                                    v-model="form.family.selected"
                                    :options="form.family.options"
                                />
                            </b-col>
                        </b-row>
                        <b-row>
                            <b-col cols="3">
                                <div class="pt-2">Textstorlek</div>
                            </b-col>
                            <b-col>
                                <b-form-input
                                    @change="onFontSizeChange"
                                    v-model.number="form.size"
                                    type="number"
                                    min="0"
                                    max="200"
                                />
                            </b-col>
                        </b-row>
                        <br />
                        <b-form-group label="Radavstånd">
                            <b-form-radio-group v-model="form.lineHeight">
                                <b-form-radio
                                    v-for="(
                                        lineHeight, indexOpt
                                    ) in form.lineHeights"
                                    :value="lineHeight"
                                    :key="indexOpt"
                                    >{{ lineHeight }}</b-form-radio
                                >
                            </b-form-radio-group>
                        </b-form-group>
                        <b-form-group label="Text- och bakgrundsfärg">
                            <b-form-radio-group v-model="form.colorID" stacked>
                                <b-form-radio
                                    v-for="(colors, indexOpt) in colorOptions"
                                    :value="colors.value"
                                    :key="indexOpt"
                                >
                                    <b-badge
                                        class="m-0 p-2"
                                        :style="{
                                            backgroundColor:
                                                colors.colors.background,
                                            color: colors.colors.foreground,
                                        }"
                                    >
                                        <span style="width: 100%">
                                            <strike
                                                v-if="
                                                    indexOpt == 6 &&
                                                    !form.customColors.valid
                                                "
                                                >{{ colors.text }}</strike
                                            >
                                            <span v-else>{{
                                                colors.text
                                            }}</span>
                                            <b-icon-brush
                                                style="margin-left: 30px"
                                                v-if="indexOpt == 6"
                                                @click="editColorScheme"
                                                v-b-tooltip.hover
                                                title="Tryck för att välja färgkombination"
                                            />
                                        </span>
                                    </b-badge>
                                </b-form-radio>
                            </b-form-radio-group>
                        </b-form-group>
                        Marginaler
                        <div class="float-right">
                            <b-form-checkbox
                                v-model="form.margins.linked"
                                @change="onMarginLinkedChange"
                                >Länka</b-form-checkbox
                            >
                        </div>
                        <br />Vänster
                        <b-form-input
                            @change="onMarginChange"
                            type="number"
                            min="0"
                            max="200"
                            v-model.number="form.margins.left"
                        />Topp
                        <b-form-input
                            @change="onMarginChange"
                            type="number"
                            min="0"
                            max="200"
                            v-model.number="form.margins.top"
                        />Höger
                        <b-form-input
                            @change="onMarginChange"
                            type="number"
                            min="0"
                            max="200"
                            v-model.number="form.margins.right"
                        />Botten
                        <b-form-input
                            @change="onMarginChange"
                            type="number"
                            min="0"
                            max="200"
                            v-model.number="form.margins.bottom"
                        />
                        <!-- <b-button type="submit" variant="primary">Spara</b-button> -->
                    </b-form>
                </div>
            </b-col>
            <b-col>
                <h4>Förhandsgranskning</h4>
                <div
                    class="fontPreview"
                    :style="{
                        fontFamily: this.form.family.selected,
                        fontSize: this.form.size + 'px',
                        backgroundColor: this.background,
                        color: this.foreground,
                        lineHeight: this.form.lineHeight,
                        paddingLeft: this.form.margins.left + 'px',
                        paddingTop: this.form.margins.top + 'px',
                        paddingRight: this.form.margins.right + 'px',
                        paddingBottom: this.form.margins.bottom + 'px',
                        overflow: 'hidden',
                    }"
                    v-html="lorem"
                />
            </b-col>
        </b-row>
    </div>
</template>

<script>
import EventBus from '../../eventbus'
import ColorPicker from '../modals/ColorPicker.vue'
export default {
    components: { ColorPicker },
    data() {
        return {
            lorem_full:
                '<b>Lorem</b> <i>ipsum</i> dolor <u>sit</u> amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. <p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>',
            lorem_short: '<b>Lo</b><i>rem</i> <u>ip</u>sum',
            foreground: '',
            background: '',

            form: {
                family: {
                    selected: 'Arial',
                    options: [
                        'Times New Roman',
                        'Arial',
                        'Roboto Mono',
                        'Verdana',
                    ],
                },
                size: 32,
                lineHeight: 1.25,
                lineHeights: [1.0, 1.25, 1.5, 1.75, 2],
                colorID: 0,
                margins: {
                    linked: true,
                    top: 0,
                    right: 0,
                    bottom: 0,
                    left: 0,
                },
                customColors: {
                    valid: false,
                    foreground: '#ffffff',
                    background: '#000000',
                },
                capitalizeOnNewLine: true,
            },
        }
    },
    computed: {
        lorem() {
            if (this.form.size > 150) {
                return this.lorem_short
            } else {
                return this.lorem_full
            }
        },
        colorOptions: {
            get: function () {
                return [
                    {
                        value: 0,
                        text: 'Vit text på svart bakgrund',
                        colors: {
                            foreground: '#ffffff',
                            background: '#000000',
                        },
                    },
                    {
                        value: 1,
                        text: 'Svart text på vit bakgrund',
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
                    {
                        value: 6,
                        text: 'Eget färgschema',
                        colors: {
                            foreground: '#FFFFFF',
                            background: '#000000',
                        },
                    },
                ]
            },
            set: function (colors) {
                return (this.colorOptions[6] = colors)
            },
        },
    },
    methods: {
        onFontSizeChange() {
            if (this.form.size <= 20) {
                this.form.size = 20
            } else if (this.form.size >= 200) {
                this.form.size = 200
            }
            this.setFontSettings()
        },
        onMarginLinkedChange() {},
        onMarginChange(val) {
            if (this.form.margins.linked) {
                console.log(val)
                if (!val) {
                    val = 0
                } else {
                    val = Math.max(
                        this.form.margins.bottom,
                        this.form.margins.top,
                        this.form.margins.left,
                        this.form.margins.right
                    )
                }
                this.form.margins.bottom =
                    this.form.margins.right =
                    this.form.margins.left =
                    this.form.margins.top =
                        val
            }
        },
        onChange() {
            this.setFontSettings()
        },
        onReset() {},

        setFontSettings() {
            this.$nextTick(() => {
                const selectedColor = this.form.colorID
                const fontSettings = {
                    font: {
                        margins: this.form.margins,
                        family: this.form.family.selected,
                        size: Number(this.form.size),
                        colorID: this.form.colorID,
                        customColors: this.form.customColors,
                        foreground:
                            this.colorOptions[selectedColor].colors.foreground,
                        background:
                            this.colorOptions[selectedColor].colors.background,
                        lineHeight: this.form.lineHeight,
                    },
                    behaviour: {
                        capitalizeOnNewLine: this.form.capitalizeOnNewLine,
                    },
                }
                this.foreground =
                    this.colorOptions[selectedColor].colors.foreground
                this.background =
                    this.colorOptions[selectedColor].colors.background

                this.$store.commit('setFontSettings', fontSettings)
            })
        },
        getSettings() {
            const fontSettings = this.$store.state.settings.font || this.form
            this.form.family.selected = fontSettings.family
            this.form.size = fontSettings.size
            this.form.colorID = fontSettings.colorID
            this.form.margins = fontSettings.margins || {
                left: 0,
                top: 0,
                right: 0,
                bottom: 0,
            }
            if (fontSettings.customColors.valid) {
                this.colorOptions[6].colors.foreground =
                    fontSettings.customColors.foreground
                this.colorOptions[6].colors.background =
                    fontSettings.customColors.background
                this.form.customColors = fontSettings.customColors
            }
            this.form.lineHeight = fontSettings.lineHeight
            const selected = this.form.colorID
            this.foreground = this.colorOptions[selected].colors.foreground
            this.background = this.colorOptions[selected].colors.background
            this.form.capitalizeOnNewLine = fontSettings.capitalizeOnNewLine
        },
        selectColor() {},
        editColorScheme(e) {
            console.log(e.preventDefault())
            this.$bvModal.show('colorPicker')
        },
        validateCustomColors(colors) {
            this.colorOptions[6].colors = colors
            this.form.customColors = colors
            let colorID = 0
            let found = this.colorOptions.slice(0, -1).some((c, i) => {
                if (
                    c.colors.background == colors.background &&
                    c.colors.foreground == colors.foreground
                ) {
                    console.log('found duplicate')
                    colorID = i
                    return true
                }
            })
            if (found) {
                this.form.customColors.valid = false
                this.form.colorID = colorID
                return false
            } else {
                this.form.customColors.valid = true
                this.form.colorID = 6
                return true
            }
        },
        updateCustomColors(colors) {
            if (this.validateCustomColors(colors)) {
                this.foreground = colors.foreground
                this.foreground = colors.foreground
            } else {
                this.foreground =
                    this.colorOptions[this.form.colorID].colors.foreground
                this.background =
                    this.colorOptions[this.form.colorID].colors.background
            }
            this.onReset()
        },
        addEventListeners() {
            EventBus.$on('savedColors', this.updateCustomColors)
            EventBus.$on('fontSettingsUpdated', this.getSettings)
        },
        removeEventListeners() {
            EventBus.$off('savedColors', this.updateCustomColors)
            EventBus.$off('fontSettingsUpdated', this.getSettings)
        },
    },
    mounted() {
        this.addEventListeners()
        this.getSettings()
    },
    beforeDestroy() {
        this.removeEventListeners()
    },
}
</script>
<style>
#dropdownMenuButton > button {
    width: 100%;
}

#dropdownMenuButton__BV_toggle_ {
    width: 100%;
}

#fontSettings {
    height: 85vh !important;
    overflow-y: scroll;
    overflow-x: hidden;
    padding-right: 1em;
}

.fontPreview {
    border: 2px solid black;
    height: 75vh !important;
    max-width: 100% !important;
    text-overflow: ellipsis;
    word-wrap: break-word;
    white-space: pre-wrap; /* CSS3 */
}
</style>
