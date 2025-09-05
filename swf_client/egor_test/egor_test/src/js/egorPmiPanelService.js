import appCtxService from 'js/appCtxService';
import soaSvc from 'soa/kernel/soaService';
import messagingService from 'js/messagingService';
import dms from 'soa/dataManagementService';
import cdm from 'soa/kernel/clientDataModel';

const RELATION_NAME = "Qam0QualityActionAttachment"
const GO_APP_URL = 'http://localhost:9000'

export async function getJTList() {
    const selection = appCtxService.ctx.xrtSummaryContextObject
    const selectionExt = await getObjects([selection.uid], [RELATION_NAME])
    const secondObjUids = selectionExt[0].props[RELATION_NAME].dbValues
    const objs = await getObjects(secondObjUids)
    return {
        totalFound: objs.length,
        searchResults: objs
    }
}

export async function getPMIServiceState(){
    const loadedJTs = await doRequest('GET','/v1/jts')
    const loadedPMIs = await doRequest('GET','/v1/pmis')
    appCtxService.updateCtx( 'egorPmiCtx', {
        loadedJTs:loadedJTs,
        loadedPMIs:loadedPMIs
    } );
}

async function doRequest(method, url ,body){
    const resp = await fetch(GO_APP_URL + url, {
            method: method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: body ? JSON.stringify(body) : {}
        })
        if (!resp.ok) throw Error("Network issue")
        const data = await resp.json()
    return data
}

export async function getPMIs(modelState) {
    try {
        if (!modelState.uid) return {
            totalFound: 0,
            searchResults: []
        }
        const uid = modelState.uid
        const data = await doRequest('POST', "/v1/jts/getPMIs",{jt_file_name: uid})
        return {
            totalFound: data.PMIs.length,
            searchResults: data.PMIs
        }
    } catch (err) {
        messagingService.showError(err.message);
        throw err
    }
}

export function log(d, view) {
    console.log(d, view)
}

export function updateParentState(selectionData, parentContext) {
    const old = parentContext.getValue()
    if (selectionData) {
        return parentContext.update(selectionData)
    } else {
        return parentContext.update(null)
    }
}

export function handleModelChangeFromMain(state) {
    if (state) return state.Title
    return ""
}

export function handlePMIChangeFromMain(state) {
    if (state) return state.props.name.value
    return ""
}

export function updatePMIInfoView(pmi) {
    const resp = []
    if (pmi.props) {
        const keys = Object.keys(pmi.props);
        for (let i = 0; i< keys.length ;i++) {
            resp.push(pmi.props[keys[i]])
        }
    }
    return {
        resp
    }
}

async function getObjects(uids, props) {
    if (props && props.length) {
        await dms.getProperties(uids, props);
    }
    let resp = await cdm.getObjects(uids);
    if (!resp || !resp.length) {
        await dms.loadObjects(uids);
        resp = await cdm.getObjects(uids);
    }
    return resp
}

    // return{
    //     totalFound:2,
    //     searchResults:[
    //         {
    //             "Title": "Model",
    //             "cellHeader1": "Model",
    //             "cellHeader2": "",
    //             "cellProperties": [ ],
    //             "hasThumbnail": false,
    //             "typeIconURL": ""
    //         },
    //         {
    //             "Title": "Body",
    //             "cellHeader1": "Body",
    //             "cellHeader2": "",
    //             "cellProperties": [],
    //             "hasThumbnail": false,
    //             "typeIconURL": ""
    //         }
    //     ]
    // }
    // if (!modelState.Title) return {
    //     totalFound:0,
    //     searchResults: []
    // }
    // const name = modelState.Title
    // return {
    //     totalFound:2,
    //     "searchResults": [
    //         {
    //             "type": "Country",
    //             "uid": -1,
    //             "props": {
    //                 "name": {
    //                     "type": "STRING",                        
    //                     "uiValue": `${name}`,
    //                     "value": `${name}`,                  
    //                 },
    //                 "value": {
    //                     "type": "STRING",                        
    //                     "uiValue": "20",
    //                     "value": "20",                  
    //                 }
    //             }
    //         },
    //         {
    //             "type": "Country",
    //             "uid": -2,
    //             "props": {
    //                 "name": {
    //                     "type": "STRING",                        
    //                     "uiValue": `${name}ZZZ`,
    //                     "value": `${name}ZZZ`,  
    //                 },
    //                 "value": {
    //                     "type": "STRING",                        
    //                     "uiValue": "40",
    //                     "value": "40",                  
    //                 }
    //             }
    //         }
    //     ]
    // }