package com.polaris.lesscode.project.internal.req;

import lombok.Data;

import java.util.List;

@Data
public class ChangeProjectChatMemberReq {

    private Long projectId;

    private List<Long> addUserIds;

    private List<Long> addDeptIds;

    private List<Long> delUserIds;

    private List<Long> delDeptIds;
}
