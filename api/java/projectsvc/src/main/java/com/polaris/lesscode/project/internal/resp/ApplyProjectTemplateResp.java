package com.polaris.lesscode.project.internal.resp;

import lombok.Data;

import java.util.Map;

@Data
public class ApplyProjectTemplateResp {

    private Long projectId;

    private Map<Long, Long> projectObjectTypeMap;

    private Map<Long, Long> iterationMap;

    private Map<Long, Long> statusMap;

}
