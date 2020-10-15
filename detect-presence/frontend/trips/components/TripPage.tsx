import React from "react";
import {useParams} from "react-router-dom";

export function TripPage() {
    const {id} = useParams<{id: string}>();

    return (
        <main>Trip: {id}</main>
    );
}
