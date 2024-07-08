// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {bilibili} from '../models';
import {main} from '../models';

export function AddAudioToList(arg1:string,arg2:string):Promise<void>;

export function AddCollectionToList(arg1:string,arg2:string,arg3:number,arg4:boolean):Promise<void>;

export function AddCompilationToList(arg1:string,arg2:number,arg3:number,arg4:number,arg5:boolean):Promise<void>;

export function AddVideoToList(arg1:string,arg2:string,arg3:boolean):Promise<void>;

export function Checkffmpeg():Promise<boolean>;

export function CreatVideoList():Promise<void>;

export function GetAppVersion():Promise<string>;

export function GetFavCollect(arg1:number):Promise<bilibili.Collects>;

export function GetListCount(arg1:string):Promise<number>;

export function GetUserInf():Promise<bilibili.AccountInformation>;

export function GetUsersCollect():Promise<bilibili.Collects>;

export function ListDownload(arg1:string,arg2:main.DownloadOption):Promise<void>;

export function LoadConfig():Promise<main.Config>;

export function LoadVideoList(arg1:string):Promise<main.VideoList>;

export function LoginBilibili():Promise<void>;

export function OpenFileDialog():Promise<string>;

export function QueryAudio(arg1:string):Promise<bilibili.Audio>;

export function QueryCollection(arg1:string):Promise<bilibili.FavList>;

export function QueryCompilation(arg1:number,arg2:number):Promise<bilibili.CompliationInformation>;

export function QuerySongInformation(arg1:string):Promise<bilibili.Audio>;

export function QueryVideo(arg1:string):Promise<bilibili.Video>;

export function RefreshConfig():Promise<void>;

export function SaveConfig(arg1:main.Config):Promise<void>;

export function SaveVideoList(arg1:main.VideoList,arg2:string):Promise<void>;
