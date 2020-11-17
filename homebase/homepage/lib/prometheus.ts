export type PrometheusResponse<Data> = PrometheusSuccessResponse<Data> | PrometheusErrorResponse<Data>

export interface PrometheusSuccessResponse<Data> {
    status: "success";
    data: Data;
    warnings?: string[];
}

export interface PrometheusErrorResponse<Data> {
    status: "error";
    data?: Data;
    errorType: string;
    error: string;
    warnings?: string[];
}

export type PrometheusLabelMap = Record<string, string>;
export type PrometheusScalarResult = [number, string];

export interface PrometheusInstantVectorResult {
    metric: PrometheusLabelMap;
    value: PrometheusScalarResult;
}

export interface PrometheusInstantVectorResults {
    resultType: "vector";
    result: PrometheusInstantVectorResult[];
}

export interface PrometheusAlert {
    activeAt?: string;
    annotations: Record<string, string>;
    labels: PrometheusLabelMap;
    state: string; // TODO sum type of values
    value: string;
}

export interface PrometheusAlerts {
    alerts: PrometheusAlert[];
}

export class PrometheusClient {
    constructor(private baseURL: string) {
    }

    async query(query: string): Promise<PrometheusInstantVectorResult[]> {
        const params = new URLSearchParams();
        params.set("query", query);
        const response = await fetch(`${this.baseURL}/api/v1/query`, {
            method: "POST",
            body: params.toString(),
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            credentials: "include",
        });
        const body: PrometheusResponse<PrometheusInstantVectorResults> = await response.json();
        if (body.status === "error") {
            throw new Error(`Prometheus error (${body.errorType}): ${body.error}`);
        }
        return body.data.result;
    }

    async alerts(): Promise<PrometheusAlert[]> {
        const response = await fetch(`${this.baseURL}/api/v1/alerts`, {
            credentials: "include",
        });
        const body: PrometheusResponse<PrometheusAlerts> = await response.json();
        if (body.status === "error") {
            throw new Error(`Prometheus error (${body.errorType}): ${body.error}`);
        }
        return body.data.alerts;
    }
}

export const client = new PrometheusClient("");
