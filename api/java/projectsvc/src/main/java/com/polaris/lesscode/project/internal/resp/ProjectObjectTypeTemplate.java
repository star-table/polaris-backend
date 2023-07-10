package com.polaris.lesscode.project.internal.resp;

import lombok.Data;

import java.util.List;

@Data
public class ProjectObjectTypeTemplate {

    private Long id;

    private String name;

    private List<ProjectStatus> status;

    private String langCode;

    private Integer sort;

    private Integer objectType;

    private Long processId;
}