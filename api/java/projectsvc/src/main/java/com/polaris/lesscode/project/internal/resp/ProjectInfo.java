package com.polaris.lesscode.project.internal.resp;

import lombok.Data;

@Data
public class ProjectInfo {
    private Long appId;

    private Long orgId;

    private Long status;

    private String name;

    private Integer publicStatus;

    private Integer publishStatus;

    private Long projectTypeId;

    private Long owner;

    private Integer resourceId;

    private Integer isFiling;

    private String remark;
}
