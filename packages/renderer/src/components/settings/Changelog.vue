<template>
  <!--<b-tabs>
    <b-tab title="Uppdateringslogg">
    -->
<div>
      <p>
        <br />
        Här listas uppdateringar, ändringar och buggfixar som gjorts i programmet
        eller distanstolkningstjänsten.
      </p>
      <b-row style="max-height: 75vh; overflow-y: scroll">
        <b-col cols="8">
          <b-card-group v-for="d in changelog" v-bind:key="'updates_on_' + d.date">
            <b-card :header="d.date | formatChangeLogDate">
              <b-list-group v-for="c in d.changes" v-bind:key="'update_' + c.id">
                <span v-if="c.type == 'version'">
                  <h3>{{ c.description }}</h3>
                </span>
                <b-list-group-item v-else>
                  <span v-if="c.type == 'fix'">
                    <b>Buggfix: </b> 
                  </span>
                  <span v-html="c.description"></span>
                  <span class="float-right" v-if="c.type != 'minor' && c.action != null">
                    <small>Åtgärd: {{ c.action }}</small>
                  </span>
                </b-list-group-item>
              </b-list-group>
            </b-card>
          </b-card-group>
        </b-col>
      </b-row>
    <!--</b-tab>
    <b-tab title="Version">
      <b-overlay :show="loading">
        <b-jumbotron>
          Du använder STTylus version:
          <code>{}</code>
          <br />Den senaste tillgängliga versionen av programmet är:
          <code>{}</code>
          <br />
        </b-jumbotron>
      </b-overlay>
    </b-tab>
  </b-tabs>
  -->
</div>
</template>

<script>
import axios from 'axios';
import EventBus from "../../eventbus";
export default {
  data() {
    return {
      loading: true,
      error: {},
      latest: new Date("2022, 12, 14"),
      changelog: [
        {
          date: new Date("2022, 12, 14"),
          changes: [
            {
              id: 0,
              type: "version",
              description: "Version 1.1.1",
            },
            {
              id: 1,
              type: "update",
              description: "Lös ut förkortningar med bortvänd parentes )",
            },
            {
              id: 2,
              type: "update",
              description: "Inaktiverar förhandsgranskning av manuskript då den inte beter sig korreket",
            },
            {
              id: 3,
              type: "update",
              description: "Ta bort föregående ord i manuskript med pil vänster",
            },
            {
              id: 4,
              type: "minor",
              description: "<i>Påminn inte</i> i stödtolkningsrutan döljer förkortningen direkt.",
            },
            {
              id: 5,
              type: "update",
              description: "Knapp för att nollställa missade förkortningar.",
            },
            {
              id: 6,
              type: "update",
              description: "Export av ämneslistor i STTylus fungerar",
            },
          ],
        },
        {
          date: new Date("2022, 12, 7"),
          changes: [
            {
              id: 0,
              type: "version",
              description: "Version 1.1.0",
            },
            {
              id: 1,
              type: "update",
              description: "Förkortningslistor cachas lokalt i webbläsaren för att skrivläget ska vara så responsivt som det bara går.",
            },
            {
              id: 2,
              type: "minor",
              description: "Tillåt ej att scrolla horisontellt i skrivläget",
            },
            {
              id: 3,
              type: "minor",
              description: "Rensa/nollställ stödtolkningsrutan när den stängs",
            },
            {
              id: 4,
              type: "update",
              description: "Åtgärdat en allvarlig bugg ('Concurrent write') som orsakat två krascher i distanstolkningsservern senaste veckorna."
            },
          ],
        },
        {
          date: new Date("2022, 11, 24"),
          changes: [ 
            {
              id: 0,
              type: "update",
              description: "Egennamn följt av kolon i början av en rad triggar stor bokstav på nästa nedslag.",
            },
            {
              id: 1,
              type: "update",
              description: "QR-kod finns som alternativ till länk.<br /> Tolkanvändarvyn kan visa QR-kod för aktuell tolkning.",
            },
          ]
        },
        {
          date: new Date("2022, 11, 23"),
          changes: [ 
            {
              id: 0,
              type: "update",
              description: "Tolkanvändarvyn försöker hindra skärmen från att gå i viloläge eller låsas."
            },
            {
              id: 1,
              type: "update",
              description: "QR-kod finns som alternativ till länk.<br /> Tolkanvändarvyn kan visa QR-kod för aktuell tolkning.",
            },
          ]
        },
        {
          date: new Date("2022, 11, 16"),
          changes: [ 
            {
              id: 0,
              type: "update",
              description: "<h2>Zoom-funktionaliteten är släckt tills vidare</h2>"
            },
            {
              id: 1,
              type: "update",
              description: "Meddela kollegan i distanstolkningen att det är dags att byta med <kbd>F3</kbd>",
            },
            {
              id: 2, 
              type:"minor",
              description: "Bättre meddelanden kring distanstolkning:<ul><li>Tolkanvändare kan ansluta i förväg, men informeras att bokningen inte börjat än</li><li>Korrekta meddelanden vid upp- och nedkoppling från tolkar och tolkanvändare.</li></ul>"
            },
            {
              id: 3,
              type: "minor",
              description: "Flyttat infomeddelanden till övre högra hörnet",
            },
            {
              id: 4,
              type: "minor",
              description: "Flyttat infomeddelanden till övre högra hörnet",
            },
            {
              id: 5,
              type:"update",
              description: "Byt standardlista med <kbd>CTRL+1...5</kbd>"
            }
          ]
        },
        {
          date: new Date("2022, 11, 10"),
          changes: [ 
            {
              id: 0,
              type: "update",
              description: "Exportera/importera mellan STTylus-konton och mellan online/offline-versionerna.",
            },
            {
              id: 1, 
              type:"update",
              description: "Optimerat förkortningsvyn kraftigt. 'Fuzzy search' i förkortningar, man behöver alltså inte söka på hela frasen eller stava rätt för att hitta matchande förkortningar."
            },
            {
              id: 2,
              type: "fix",
              description: "Alla dialogrutor stänger om man trycker på escape",
            },
            {
              id: 3,
              type: "fix",
              description: "<kbd>CTRL+HÖGERPIL</kbd> stannar vid sista positionen i texten.",
            }
          ]
        },
        {
          date: new Date("2022, 11, 4"),
          changes: [ 
            {
              id: 0,
              type: "fix",
              description: "När förkortningar matas ut på sista raden/nästa rad så följer textmarkören med.",
            },
            {
              id: 1, 
              type:"fix",
              description: "Åtgärdat en bugg där det blir stor bokstav på konstiga ställen om man backar förbi ett stort skiljetecken med <kbd>BACKSPACE</kbd> eller <kbd>CTRL+BACKSPACE</kbd>.",
            },
            {
              id: 2,
              type:"update",
              description: "Offline-versionen av STTylus använder AppData så att användarinställningar sparas mellan uppdateringarna", 
            },
            {
              id: 3,
              type: "update",
              description:"<kbd>F11</kbd> växlar till och från fullskärm i offline-versionen.",
            }
          ]
        },
        {
          date: new Date("2022, 08, 26"),
          changes: [ 
            {
              id: 0,
              type: "update",
              description: "Förhandsgranskning av manuskript.",
            },
            {
              id: 1, 
              type:"fix",
              description: "Åtgärdat en bugg där talstreck i början av en rad förvandlas till en punktlista."
            }
          ]
        },
        {
          date: new Date("2022, 06, 27"),
          changes: [ 
            {
              id: 0,
              type: "update",
              description: "<h2>STTylus finns nu som desktop- och offlineversion.</h2>",
            }
          ]
        },
        {
          date: new Date("2022, 04, 01"),
          changes: [
            {
              id: 0,
              type: "fix",
              description:
                "Åtgärdat en bugg som gjorde att man inte kunde skriva ? eller /.",
            },
            {
              id: 1,
              type: "update",
              description:
                "Skapat funktioner för exportering av listor till ProType och TextOnTop.",
              action: "Återkoppla gärna om TextOnTop, jag har inte det programmet själv"
            },
            {
              id: 2,
              type: "update",
              description:
                "Skapat funktion för importering av IllumiType-exporter.",
            },
            {
              id: 3,
              type: "update",
              description: "Fälten i Stödtolkningsrutan kan visas/döljas och markeras med kortkommandon."
            },
            {
              id: 4,
              type: "fix",
              description: "Delad lista funktionen har bättre och säkrare felhantering."
            }
          ],
        },
        {
          date: new Date("2022/03/24"),
          changes: [
            {
              id: 0,
              type: "fix",
              description:
                "Manuskript-redigering upptäcker osparade ändringar och stänger öppnade dokument korrekt.",
            },
            {
              id: 1,
              type: "minor",
              description:
                "<kbd>ENTER</kbd> ansluter till kollegas distanstolkning istället för att man ska behöva klicka med musen.",
            },
            {
              id: 2,
              type: "update",
              description:
                "Stänger man ner webbläsaren eller fliken STTylus körs i så kopplas man ner från eventuell distanstolkning.",
            }
          ],
        },
        {
          date: new Date("2022/03/12"),
          changes: [
            {
              id: 0,
              type: "update",
              description:
                "STTylus tar diverse genvägar för att använda mindre bandbredd.<br />Förhoppningen är att det resulterar i att programmet fungerar bättre när/om man råkar ha svajig internetuppkoppling.",
            },
          ],
        },
        {
          date: new Date("2022/03/11"),
          changes: [
            {
              id: 0,
              type: "fix",
              description:
                "Manuskript laddas in korrekt igen.<br/>En tidigare uppdatering ställde till det så att man behövde ladda om hela sidan för att manuskripten skulle läsas in.",
            },
            {
              id: 1,
              type: "minor",
              description:
                "Stödtolkningsrutan markerar 'Slå upp förkortning' om delad lista inte används. Annars markeras 'Lägg till delad förkortning'",
            },
          ],
        },
        {
          date: new Date("2022/03/10"),
          changes: [
            {
              id: 0,
              type: "minor",
              description:
                "<kbd>SHIFT+F2</kbd> öppnar 'Skapa förkortning' men ger frasen små bokstäver.",
            },
            {
              id: 1,
              type: "minor",
              description: "Inställning för byte av typsnitt.",
            },
            {
              id: 2,
              type: "minor",
              description: "Tolkanvändare kan byta radavstånd och typsnitt.",
            },
            {
              id: 3,
              type: "update",
              description:
                "Lagt till en högerspalt som visar anslutna användare/tolkar under distanstolkning.<br/>Öppna och stäng den med <kbd>F3</kbd> eller genom att klicka på 'Visa ... anslutna' i toppmenyn.",
              action: "Testa och återkom med feedback",
            },
          ],
        },
        {
          date: new Date("2022/03/09"),
          changes: [
            {
              id: 0,
              type: "minor",
              description:
                "Förkortningar matas fram när man scrollar i listan. <br />Scrollbar lista över förkortningslistor.",
            },
            {
              id: 1,
              type: "minor",
              description:
                "Visar antalet anslutna till din distanstolkning i toppmenyn.",
            },
            {
              id: 2,
              type: "update",
              description:
                "STTylus laddar in inställningssidan i bakgrunden vilket gör att det går mycket snabbare att växla fram och tillbaka till textläget.",
              action: "Testa och återkoppla om fel uppstår.",
            },
            {
              id: 3,
              type: "minor",
              description:
                "Meddelandena <em>Tolk kopplade upp/ner</em> ska endast gälla när en tolk faktiskt kopplar upp och ner som en följd av den större uppdateringen ovan.",
            },
          ],
        },
        {
          date: new Date("2022/03/07"),
          changes: [
            {
              id: 0,
              type: "minor",
              description:
                "Notifikationer visas längre tid och med större text.",
            },
            {
              id: 1,
              type: "minor",
              description: "Större text i 'Lägg till förkortning'.",
            },
            {
              id: 2,
              type: "minor",
              description:
                "Förkortningar kan slås ut på tre sätt:<br/><ol><li><em>smsm</em> blir <em>så mycket som möjligt</em></li><li><em>Smsm</em> blir <em>Så mycket som möjligt</em></li><li><em>SMSM</em> blir <em>SÅ MYCKET SOM MÖJLIGT</em></li></ol>",
            },
          ],
        },
        {
          date: new Date("2022/02/10"),
          changes: [
            {
              id: 0,
              type: "minor",
              description:
                "'Ignorera alla' i stödtolkningsrutan ska nu fungera.<br />Om stödtolksfunktionen är långsam på att reagera så kan det bero på att listan med förslag har blivit lång. Detta problem i sig ska lösas på sikt.",
            },
          ],
        },
        {
          date: new Date("2022/02/09"),
          changes: [
            {
              id: 0,
              type: "minor",
              description:
                "Förkortningar löses ut efter citat-tecken (' och \").",
            },
            {
              id: 1,
              type: "minor",
              description:
                "STTylus varnar när programmet inte kan skicka text till Zoom CC.",
            },
          ],
        },
        {
          date: new Date("2022/01/25"),
          changes: [
            {
              id: 0,
              type: "update",
              description:
                "Lagt till möjligheten att redigera förkortningslistor (namn och typ).",
              action: "Testa och återkoppla om fel uppstår",
            },
            {
              id: 1,
              type: "minor",
              description: "Ändrat färgschemat",
            },
          ],
        },
        {
          date: new Date("2022/01/13"),
          changes: [
            {
              id: 0,
              type: "update",
              description:
                "Lagt till möjligheten att importera TextOnTop-listor.",
              action: "Testa och återkoppla om fel uppstår",
            },
          ],
        },
        {
          date: new Date("2021/12/13"),
          changes: [
            {
              id: 0,
              type: "update",
              description:
                "Lagt till stödtolkningsfunktionen 'Delad lista'.<br>Vid distanstolkning kan du skapa en delad lista som både du och din kollega får tillgång till.<br />Öppna stödtolkningsfunktionen och tryck <kbd>Enter</kbd> i rutan uppe i höger för att skapa den. Sifferkoden delar du med din kollega.<br />Förkortningar som läggs till i det här läget kan användas av båda skrivtolkarna.<br />Efter distanstolkningen kan du välja att importera och på så sätt spara förkortningar från listan.",
              action: "Testa och återkoppla om fel uppstår",
            },
            {
              id: 1,
              type: "minor",
              description:
                "Importering av förkortningslistor hanterar felaktiga format och konflikter på ett bättre sätt. Det går även att ladda in i befintliga listor om konflikterna blir lösta.",
            },
          ],
        },
        {
          date: new Date("2021/12/02"),
          changes: [
            {
              id: 0,
              type: "fix",
              description:
                'Åtgärdat en bugg där frammatning av manuskript ignorerade upprepade skiljetecken i slutet på en rad.<br />Exempel: <i>"Det här är ett citat på en rad!"</i> eller <i>En mening avslutas med trippelpunkt/ellips...</i>.<br />I dessa fall skulle bara det första skiljetecknet tas med, så programmet skulle ha matat fram <i>"det här är ett citat<b>!</b></i> och <i>trippelpunkt/ellipsis.</i>',
              action: "Testa och återkoppla om fel uppstår",
            },
            {
              id: 1,
              type: "minor",
              description:
                "Frammatning av manuskript stängs av oavsett hur slutet av dokumentet nåtts.",
              action: "Testa och återkoppla om fel uppstår",
            },
            {
              id: 2,
              type: "minor",
              description:
                "Lösenord kan nollställas på så sätt att tolken själv väljer sitt nya lösenord.",
            },
            {
              id: 3,
              type: "minor",
              description:
                "Visar hjälptext för hur man byter lista i <b>Skapa förkortning</b> om man har mer än en lista markerad.<br /><kbd>F2</kbd> både öppnar rutan och växlar mellan listor.",
            },
          ],
        },
        {
          date: new Date("2021/11/30"),
          changes: [
            {
              id: 0,
              type: "fix",
              description:
                'Systemet rättar till förkortningar med stora bokstäver.<br />Tidigare så bråkar systemet om man lägger in t.ex. <i>Hoas</i> för hoppas.<br />Förkortningar med inledande stor bokstav i förkortningen kan inte lösas ut i skrift och förkortningslistan tolkar <i>Hoas</i> som en annan förkortning än <i>hoas</i>.<br />Detta gäller både om "Skapa förkortning" eller <kbd>F2</kbd> används, men också vid importering av förkortningar',
              action: "Testa och återkoppla om fel uppstår",
            },
          ],
        },
        {
          date: new Date("2021/11/19"),
          changes: [
            {
              id: 0,
              type: "update",
              description:
                "Manuskriptfunktionen ska vara i rullning och även fungera tillsammans med Zoom.<br />Mata fram manuskript på samma rad som dess förkortning med <kbd>Space</kbd> eller på ny rad med <kbd>Enter</kbd>.",
              action: "Testa och återkoppla om fel uppstår",
            },
            {
              id: 1,
              type: "fix",
              description:
                "Problem där Zoom inte tog emot text när man skrev på annan rad än den längst ner.<br />Gör man ny rad när man har en tom rad under markören så hoppar den längst ner.",
              action: "Testa och återkoppla om fel uppstår",
            },
            {
              id: 2,
              type: "minor",
              description:
                "Radavstånd mindre än 1 är borrtaget då t.ex. Å, Ä och Ö visades konstigt med den inställningen.",
            },
          ],
        },
      ],
    };
  },
  computed: {
    version() {
      return {
        version: "0.0.0",
        frontend: "0.0.0",
        backend: "0.0.0",
        latest: {},
      }
    }
  },
  methods: {
    getFrontendVer() {

    },
    getBackendVer() {

    },
    getConsumerVer() {

    },
    getManifest() {
      this.loading = true
      axios.get("http://localhost:3000/").then(resp => {
        this.latest = resp.data
        this.loading = false
      }).catch(err => {
        this.loading = false
        this.error = err
      })
    }
  },
  mounted() {
    this.getManifest()
    if (this.$store.state.lastLogin > this.latest) {
      console.log("number of changes:", this.changelog[0].changes.length);
      setTimeout(() => {
        EventBus.$emit("newChangeLog", this.changelog[0].changes.length);
      }, 125);
    }
  },
};
</script>
