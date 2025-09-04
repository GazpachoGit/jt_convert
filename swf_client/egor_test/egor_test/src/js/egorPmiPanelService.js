import appCtxService from 'js/appCtxService';
import soaSvc from 'soa/kernel/soaService';
import messagingService from 'js/messagingService';
import dms from 'soa/dataManagementService';
import cdm from 'soa/kernel/clientDataModel';


export async function getJTList(){
    // const objs = await getObjects(["AgDAA0_CpqqTBB"])
    // return {
    //     totalFound:1,
    //     searchResults:objs
    // }
    return{
        totalFound:2,
        searchResults:[
            {
                "Title": "Model",
                "cellHeader1": "Model",
                "cellHeader2": "",
                "cellProperties": [ ],
                "hasThumbnail": false,
                "typeIconURL": ""
            },
            {
                "Title": "Body",
                "cellHeader1": "Body",
                "cellHeader2": "",
                "cellProperties": [],
                "hasThumbnail": false,
                "typeIconURL": ""
            }
        ]
    }
}

async function getObjects(uids, props){
    if(props && props.length){
        await dms.getProperties(uids, props);
    }
    let resp = await cdm.getObjects(uids);
    if(!resp || !resp.length){
        await dms.loadObjects(uids);
        resp  = await cdm.getObjects(uids);
    }
    return resp
}

export function getPMIs(modelState){
    if (!modelState.Title) return {
        totalFound:0,
        searchResults: []
    }
    const name = modelState.Title
    return {
        totalFound:2,
        "searchResults": [
            {
                "type": "Country",
                "uid": -1,
                "props": {
                    "name": {
                        "type": "STRING",                        
                        "uiValue": `${name}`,
                        "value": `${name}`,                  
                    },
                    "value": {
                        "type": "STRING",                        
                        "uiValue": "20",
                        "value": "20",                  
                    }
                }
            },
            {
                "type": "Country",
                "uid": -2,
                "props": {
                    "name": {
                        "type": "STRING",                        
                        "uiValue": `${name}ZZZ`,
                        "value": `${name}ZZZ`,  
                    },
                    "value": {
                        "type": "STRING",                        
                        "uiValue": "40",
                        "value": "40",                  
                    }
                }
            }
        ]
    }
}

export function log(d, view){
    console.log(d,view)
}

export function updateParentState(selectionData, parentContext){
    const old = parentContext.getValue()
    if ( selectionData ) {
        return parentContext.update( selectionData )
    } else {
        return parentContext.update( null )
    }
}

export function handleModelChangeFromMain(state){
    if (state) return state.Title
    return ""
}

export function handlePMIChangeFromMain(state){
    if (state) return  state.props.name.value 
    return ""
}

export function updatePMIInfoView(pmi){
    let name, value = ""
    if (pmi.props){
        name = pmi.props.name.uiValue,
        value = pmi.props.value.uiValue
    }
    return {
        name,
        value
    }
}