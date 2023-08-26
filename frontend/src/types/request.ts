import { AxiosResponse } from 'axios';

export type PromiseResp<T> = Promise<AxiosResponse<T>>
export type Resp<T> = PromiseResp<T>;
