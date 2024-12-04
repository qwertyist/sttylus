
        <!--<b-row v-if="collab">
            <b-col cols="6">
                <h4>Lägg till delad förkortning</h4>
            </b-col>
            <b-col>
                <b-form @submit.prevent="onSubmitSharedList" autocomplete="off">
                    <b-row>
                        <b-col cols="2"></b-col>
                        <b-col cols="2">
                            <b-button
                                size="sm"
                                v-if="sharingList"
                                @click.prevent="onLeaveSharedList"
                                >Lämna</b-button
                            >
                        </b-col>
                        <b-col>
                            <b-form-select
                                size="sm"
                                :disabled="form.sharedList != ''"
                                v-model="selectedBaseList"
                            >
                                <option value="null">Ny lista</option>
                                <option
                                    v-for="list in selectedLists"
                                    :value="list.id"
                                    :key="list.id"
                                    :disabled="form.sharedList != ''"
                                >
                                    {{ list.name }}
                                </option>
                            </b-form-select>
                            <b-form-input
                                :disabled="sharingList"
                                v-model="form.sharedList"
                                id="list"
                                size="sm"
                                placeholder="Skapa/gå med i delad lista"
                                ref="sharedListInput"
                                spellcheck="false"
                            />
                        </b-col>
                    </b-row>
                </b-form>
            </b-col>
        </b-row>
        <div v-if="collab">
            <b-form
                @submit.prevent="onSubmitSharedAbb"
                autocomplete="off"
                class="d-flex justify-content-between align-items-center"
            >
                <b-form-input
                    v-model="form.sharedAbb"
                    :disabled="!sharingList"
                    id="abb"
                    placeholder="Förkortning"
                    autofocus
                    style="width: 20%"
                    ref="sharedAbbInput"
                />
                <b-form-input
                    v-model="form.sharedWord"
                    :disabled="!sharingList"
                    id="word"
                    placeholder="Ord"
                    style="width: 60%"
                    ref="sharedWordInput"
                />
                <b-form-group>
                    <b-badge class="mt-3">Fokusera med CTRL+1</b-badge>
                </b-form-group>
                <input type="submit" style="display: none" />
            </b-form>
            <hr />
        </div>
        <div v-if="sharedAbbs">
            <div class="d-flex justify-content-between align-items-center">
                <span
                    style="
                        font-size: 1.5rem;
                        margin-bottom: 0.5rem;
                        font-weight: 500;
                        line-height: 1.2;
                    "
                    >Delade förkortningar</span
                >
                <span>
                    <b-badge class="float-right">Visa/dölj med CTRL+2</b-badge>
                    <div v-if="show.sharedAbbs">
                        <b-badge>Växla bokstav med CTRL+Bokstav/Noll</b-badge>
                    </div>
                    <br />
                </span>
                <b-col v-if="!collab && selectedBaseList">
                    <small class="mt-3">
                        Du är nedkopplad, men kan fortfarande spara delade
                        förkortningar.
                    </small>
                    <b-button
                        class="float-right"
                        size="sm"
                        v-if="sharingList"
                        variant="primary"
                        @click.prevent="onSaveSharedAbbs"
                        >Spara</b-button
                    >
                </b-col>
            </div>
            <div v-if="show.sharedAbbs">
                <span v-for="letter in letters" :key="letter" class="p-2 mb-3">
                    <b v-if="letter == currentLetter">{{ letter }}</b>
                    <a v-else href="#" @click="setCurrentLetter(letter)">{{
                        letter
                    }}</a>
                </span>
                <b-list-group ref="sharedAbbsGroup">
                    <b-list-group-item
                        class="p-1 pl-3"
                        v-for="(abb, i) in sharedAbbsByLetter"
                        :key="'shared_' + i"
                    >
                        {{ abb.abb }} - {{ abb.word }}
                        <b-button
                            size="sm"
                            variant="danger"
                            @click.prevent="onRemoveSharedAbb(abb.abb)"
                            class="float-right"
                        >
                            <b-icon icon="trash" />
                        </b-button>
                    </b-list-group-item>
                </b-list-group>
            </div>
            <div v-else>
                <b-badge>...</b-badge>
            </div>
            <hr />
        </div>
        -->
