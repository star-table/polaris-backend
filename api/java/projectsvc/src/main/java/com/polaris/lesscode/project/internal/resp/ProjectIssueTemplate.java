package com.polaris.lesscode.project.internal.resp;

import lombok.Data;

@Data
public class ProjectIssueTemplate {

    private String title;

    private Long projectObjectTypeId;

    private Long statusId;

    private Long iterationId;

}
