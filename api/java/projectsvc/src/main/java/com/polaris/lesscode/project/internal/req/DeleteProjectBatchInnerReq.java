package com.polaris.lesscode.project.internal.req;

import lombok.Data;

import java.util.List;

@Data
public class DeleteProjectBatchInnerReq {
    private List<Long> ProjectIds;
}
