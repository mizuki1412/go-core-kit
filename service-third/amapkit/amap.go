package amapkit

/** lon, lat */

//public Arr<String,Object> geo(String cityCode, String address){
//Arr<String,Object> map = new HashMap<>();
//map.put("key",KEY);
//map.put("address",address);
//map.put("city",cityCode+"00");
//ResponseEntity<String> responseEntity = restTemplate.getForEntity("https://restapi.amap.com/v3/geocode/geo?key={key}&address={address}&city={city}",String.class,map);
//Arr<String,Object> data = JsonUtil.toMap(responseEntity.getBody());
//if(data==null || "0".equals(String.valueOf(data.getOrDefault("status",0)))) return null;
//List<Arr<String,Object>> list = (List<Arr<String, Object>>) data.getOrDefault("geocodes",new ArrayList<>());
//if(list.size()>0){
//Arr<String,Object> e = list.get(0);
//String location = (String)e.get("location");
//if(location!=null){
//location = location.substring(0,location.length()-1);
//String[] locations = location.split(",");
//BigDecimal lon = BigDecimal.valueOf(Double.parseDouble(locations[0]));
//BigDecimal lat = BigDecimal.valueOf(Double.parseDouble(locations[1]));
//Arr<String,Object> ret = new HashMap<>();
//ret.put("lon",lon);
//ret.put("lat",lat);
//return ret;
//}
//}
//return null;
//}
//
///** pca, address */
//public Arr<String,Object> regeo(BigDecimal lon, BigDecimal lat){
//Arr<String,Object> map = new HashMap<>();
//map.put("key",KEY);
//map.put("location",lon.toString()+","+lat.toString());
//ResponseEntity<String> responseEntity = restTemplate.getForEntity("https://restapi.amap.com/v3/geocode/regeo?key={key}&location={location}",String.class,map);
//Arr<String,Object> data = JsonUtil.toMap(responseEntity.getBody());
//if(data==null || "0".equals(String.valueOf(data.getOrDefault("status",0)))) return null;
//Arr<String,Object> regeocode = (Arr<String, Object>) data.get("regeocode");
//if(regeocode==null) return null;
//String address = (String) regeocode.get("formatted_address");
//if(address==null) return null;
//Arr<String,Object> addressComponent = (Arr<String, Object>) regeocode.get("addressComponent");
//if(addressComponent==null) return null;
//Arr<String,Object> ret = new HashMap<>();
//ret.put("provinceName",addressComponent.get("province"));
//ret.put("cityName",addressComponent.get("city"));
//ret.put("areaName",addressComponent.get("district"));
//ret.put("address",address.replace(ret.get("provinceName").toString()+ret.get("cityName")+ret.get("areaName"),""));
//String adcode = (String) addressComponent.get("adcode");
//if (adcode==null) return null;
//ret.put("provinceCode",adcode.substring(0,2));
//ret.put("cityCode",adcode.substring(0,4));
//ret.put("areaCode",adcode);
//return ret;
//}
//
//@Autowired
//private ProjectConfig projectConfig;
//
//public List<Arr<String,Object>> weather() throws RestMainException {
//if(projectConfig.getWeatherCity()==null) throw new RestMainException("天气城市未设置");
//Arr<String,Object> map = new HashMap<>();
//map.put("key",KEY);
//map.put("city",projectConfig.getWeatherCity());
//map.put("extensions","all");
//ResponseEntity<String> responseEntity = restTemplate.getForEntity("https://restapi.amap.com/v3/weather/weatherInfo?key={key}&city={city}&extensions={extensions}",String.class,map);
//Arr<String,Object> data = JsonUtil.toMap(responseEntity.getBody());
//if(data==null || "0".equals(String.valueOf(data.getOrDefault("status",0)))) return null;
//List<Arr<String,Object>> list = (List<Arr<String, Object>>) data.getOrDefault("forecasts",new ArrayList<>());
//if(list.size()>0){
//return (List<Arr<String,Object>>)list.get(0).getOrDefault("casts",new ArrayList<>());
//}
//return new ArrayList<>();
//}
