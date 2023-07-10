package com.polaris.lesscode.project.internal.resp;

import lombok.Data;

@Data
public class IterationStatus {

    private String startTime;

    private String endTime;

    private String planStartTime;

    private String planEndTime;

    private Integer statusType;
}
