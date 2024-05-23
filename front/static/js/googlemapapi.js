/**
 * @license
 * Copyright 2024 Google LLC. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
async function init() {
    await customElements.whenDefined('gmp-map');

    const map = document.querySelector("gmp-map");
    const marker = document.getElementById("marker");
    const strictBoundsInputElement = document.getElementById("use-strict-bounds");
    const placePicker = document.getElementById("place-picker");
    const infowindowContent = document.getElementById("infowindow-content");
    const infowindow = new google.maps.InfoWindow();

    map.innerMap.setOptions({mapTypeControl: false});
    infowindow.setContent(infowindowContent);

    placePicker.addEventListener('gmpx-placechange', () => {
        const place = placePicker.value;

        if (!place.location) {
            window.alert(
                "No details available for input: '" + place.name + "'"
            );
            infowindow.close();
            marker.position = null;
            return;
        }

        if (place.viewport) {
            map.innerMap.fitBounds(place.viewport);
        } else {
            map.center = place.location;
            map.zoom = 17;
        }

        marker.position = place.location;
        infowindowContent.children["place-name"].textContent = place.displayName;
        infowindowContent.children["place-address"].textContent = place.formattedAddress;
        infowindow.open(map.innerMap, marker);
    });

    // Sets a listener on a radio button to change the filter type on the place picker
    function setupClickListener(id, type) {
        const radioButton = document.getElementById(id);
        radioButton.addEventListener("click", () => {
            placePicker.type = type;
        });
    }
    setupClickListener("changetype-all", "");
    setupClickListener("changetype-address", "address");
    setupClickListener("changetype-establishment", "establishment");
    setupClickListener("changetype-geocode", "geocode");
    setupClickListener("changetype-cities", "(cities)");
    setupClickListener("changetype-regions", "(regions)");

    strictBoundsInputElement.addEventListener("change", () => {
        placePicker.strictBounds = strictBoundsInputElement.checked;
    });
}

document.addEventListener('DOMContentLoaded', init);