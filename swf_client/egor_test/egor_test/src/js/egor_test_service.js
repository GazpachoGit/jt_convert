import appCtx from 'js/appCtxService';

export function test(){
    const xrtContext = appCtx.getCtx( 'ActiveWorkspace:xrtContext' );
    console.log("xrtContext: ")
    console.log(xrtContext)
}