package com.polaris.lesscode.msgsvc.internal.req;

import lombok.Data;

import java.util.List;

@Data
public class GlobalRefresh {

    private long objectId;

    private Object objectValue;

    private Object newObjectValue;

    private List<Long> childrenIssueIds;
}

