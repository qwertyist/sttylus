<template>
    <div id="hotkeys">
        <b-row style="height: 85vh !important; overflow-y: auto">
            <b-col cols="4">
                <b-list-group-item
                    class="list-group col-lg-12"
                    v-for="(hotkey, index) in functionkeys"
                    :key="index"
                >
                    <kbd>{{ hotkey.whichKey }}</kbd>
                    &nbsp;&nbsp;<span v-html="hotkey.command" />
                </b-list-group-item>
            </b-col>
            <b-col cols="7">
                <b-list-group-item
                    class="list-group col-lg-12"
                    v-for="(hotkey, index) in hotkeys"
                    :key="index"
                >
                    <kbd>{{ hotkey.whichKey }}</kbd>
                    {{ hotkey.command }}
                </b-list-group-item>
            </b-col>
        </b-row>
    </div>
</template>

<script>
export default {
    data() {
        return {
            functionkeys: [
                { whichKey: 'F1', command: 'Visa stödtolkningsfunktionen' },
                {
                    whichKey: 'F2',
                    command:
                  'Lägg till en förkortning<br />I skrivläget så väljs ordet vid/bakom textmarkören<br />Markera flera ord med <kbd>CTRL+SHIFT ←/→</kbd><br /><kbd>SHIFT+F2</kbd> ger den valda frasen gemener',
                },
                { whichKey: 'F2', command: 'Växla mellan förkortningslistor i Lägg till-rutan' },
                { whichKey: 'F3', command: 'Signalera redo för byte' },
                { whichKey: 'F4', command: 'Rensa texten' },
                {
                    whichKey: 'F5',
                    command: 'Öppna och stäng programinställningar',
                },
                { whichKey: 'F6', command: 'Gör texten mindre' },
                { whichKey: 'F7', command: 'Gör texten större' },
                { whichKey: 'F8', command: 'Växla textfärg' },
                { whichKey: 'F10', command: 'Visa chatt/Växla fokus mellan text och chatt' },
                { whichKey: 'SHIFT+F10', command: 'Dölj chatt' },
                { whichKey: 'F11', command: 'Växla helskärm' },
            ],
            hotkeys: [
                {
                    whichKey: 'ESC',
                    command:
                        'Hoppar till slutet av dokumentet, döljer toppmenyn samt avbryter frammatning av manuskript',
                },
                { whichKey: 'CTRL+B', command: 'Växla fetstilt text' },
                { whichKey: 'CTRL+I', command: 'Växla kursiv text' },
                { whichKey: 'CTRL+U', command: 'Växla understruken text' },
                {
                    whichKey: 'CTRL+SPACE',
                    command: 'Lös ut förkortning utan mellanslag',
                },
                {
                    whichKey: 'TAB',
                    command: 'Avbryt frammatning av förkortning. Andra nedslaget infogar ett mellanslag.',
                },
                {
                    whichKey: 'SHIFT+ENTER',
                    command:
                        'Infoga ny rad med utan att versalera nästa ord'
                },
                {
                    whichKey: 'CTRL+K',
                    command:
                        'Infoga Protype-klubban vid markören'
                },
                {
                    whichKey: 'PIL HÖGER →',
                    command:
                        'Mata fram ord ur aktiverat manuskript vid sista ordet',
                },
                {
                    whichKey: 'PIL NEDÅT ↓',
                    command: 'Mata fram ord till nästa skiljetecken',
                },
                {
                    whichKey: 'PIL VÄNSTER ←',
                    command: 'Ta bort senaste ordet ur aktiverat manuskript',
                },
                {
                    whichKey: 'CTRL+SHIFT+V',
                    command: 'Klistra in (oformaterad) text i textvyn',
                },
                {
                    whichKey: 'CTRL+1...5',
                    command: 'Växla mellan standardlistor',
                },
                /*{
          whichKey: "CTRL+HÖGER PIL",
          command: "Flytta markören åt höger vid användning av manuskript",
        },
        {
          whichKey: "CTRL+VÄNSTER PIL",
          command: "Flytta markören åt vänster vid användning av manuskript",
        },*/
            ],
        }
    },
}
</script>
<style>
kbd {
    background-color: #eee;
    border-radius: 3px;
    border: 1px solid #b4b4b4;
    box-shadow:
        0 1px 1px rgba(0, 0, 0, 0.2),
        0 2px 0 0 rgba(255, 255, 255, 0.7) inset;
    color: #333;
    display: inline-block;
    font-size: 0.85em;
    font-weight: 700;
    line-height: 1;
    padding: 2px 4px;
    white-space: nowrap;
}
</style>
