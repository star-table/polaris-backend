package com.polaris.lesscode.msgsvc.internal.req;

import lombok.Data;

import java.util.List;

@Data
public class PushMsgToMqttReq {

    private long orgId;

    private long projectId;

    private long appId;

    private String action;

    private String type;

    private List<PartialRefresh> partialRefresh;

    private List<GlobalRefresh> globalRefresh;

}
