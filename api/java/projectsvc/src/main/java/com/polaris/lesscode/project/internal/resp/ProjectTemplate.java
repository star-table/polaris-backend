package com.polaris.lesscode.project.internal.resp;

import lombok.Data;

import java.util.List;
import java.util.Map;

@Data
public class ProjectTemplate {
    private Integer templateFlag;

    private List<ProjectObjectTypeTemplate> projectObjectTypeTemplates;

    private List<ProjectIterationTemplate> projectIterationTemplates;

    private List<Map<String, Object>> projectIssueTemplates;

    private ProjectInfo projectInfo;

    private ProjectDetail projectDetail;

    private TemplateInfo templateInfo;

    private Boolean isNewbieGuide;

    private Boolean needData;

    private Boolean isCreateTemplate;

    private Long fromOrgId;
}
