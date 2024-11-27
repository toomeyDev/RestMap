<!-- src/components/Map.vue -->
<template>
  <div id="map-component" style="width: 100%; height: 100vh"></div>
</template>

<script>
  import mapboxgl from "mapbox-gl";
  import axios from 'axios';

  export default {
    name: 'Map-Component',
    data () {
      return {
        map: null,
      };
    },
    mounted() {
      mapboxgl.accessToken = 'pk.eyJ1IjoibXRvb21leS1kZXYiLCJhIjoiY200MDRiOWxqMTlzODJ3cHlrejE3NjgzNyJ9.AqCgaLym-Ym9jBSrsj6aGg';
      this.map = new mapboxgl.Map({
        container: 'map-component',
        style: 'mapbox://styles/mapbox/streets-v11',
        center: [-74.0060, 40.7128], // Starting pos. (lng, lat) NYC
        zoom: 12,
      });

      // Fetch restaurant data
      axios.get('http://localhost:8080/restaurants')
        .then(response => {
          const restaurants = response.data.restaurants;
          restaurants.forEach(restaurant => {
            new mapboxgl.Marker()
              .setLngLat([restaurant.longitude, restaurant.latitude])
              .setPopup(new mapboxgl.Popup().setText(restaurant.name))
              .addTo(this.map);
          });
        })
        .catch(error => {
          console.error('Error fetching restaurant data: ', error);
        });
    },
  };
</script>

<style>
  #map-component {
    position: relative;
  }
</style>
