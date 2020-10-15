import React from 'react';

export function TripTag({tag}: {tag: string}) {
    return (
        <div className="inline-flex py-1 px-2 mx-1 rounded font-bold text-xs bg-gray-400 text-gray-800">
            {tag}
        </div>
    )
}
