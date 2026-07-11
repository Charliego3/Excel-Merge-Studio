import type { Component } from "svelte";

export interface Flow {
    steps: FlowStep[];
}

export interface FlowStep {
    name: string;
    description: string;
    path: string;
    icon: Component;
    iconColor: string;
    iconBgColor: string;
}
