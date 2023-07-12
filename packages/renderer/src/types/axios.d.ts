import { Axios } from 'axios'

declare module 'axios' {
    export interface AxiosRequestConfig {
        createAbb: boolean
        clearCache: boolean
        identifier: string
    }
}
