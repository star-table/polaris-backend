package com.polaris.lesscode.project.internal.resp;

import lombok.Data;

import java.util.List;

@Data
public class ProjectIterationTemplate {

    private Long id;

    private String name;

    private String startTime;

    private String endTime;

    private Long sort;

    private Long owner;

    private Integer statusType;

    private List<IterationStatus> iterationStatus;

}
