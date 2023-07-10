package com.polaris.lesscode.msgsvc.internal.api;


import com.polaris.lesscode.msgsvc.internal.req.PushMsgToMqttReq;
import com.polaris.lesscode.vo.Result;
import org.springframework.web.bind.annotation.*;


@RequestMapping("/api/msgsvc")
public interface MsgApi {

    @PostMapping(value = "/pushMsgToMqtt")
    Result<?> pushMsgToMqtt(@RequestBody PushMsgToMqttReq req);


}
