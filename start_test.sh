#!/bin/bash

#curl -XPUT -v michael:123456@127.0.0.1:5000/v2/service_instances/123 -d '{"organization_guid": "org-guid-here", "plan_id": "3C7BAF72-0DF4-420D-B365-B5CF09409C70","service_id":"04EB4D8F-15BF-43F2-B4DA-E7A243E21C83"}'
#sleep 3
#curl -XPUT -v michael:123456@127.0.0.1:5000/v2/service_instances/123/service_bindings/456 -d '{"app_gui": "org-guid-here", "plan_id": "3C7BAF72-0DF4-420D-B365-B5CF09409C70","service_id":"04EB4D8F-15BF-43F2-B4DA-E7A243E21C83"}'
#sleep 3
#curl -XDELETE michael:123456@127.0.0.1:5000/v2/service_instances/123?service_id=04EB4D8F-15BF-43F2-B4DA-E7A243E21C83\&plan_id=3C7BAF72-0DF4-420D-B365-B5CF09409C70
#curl -XDELETE -v michael:123456@127.0.0.1:5000/v2/service_instances/123/service_bindings/456

curl -XGET  michael:123456@127.0.0.1:5000/v2/catalog
curl -XPUT  michael:123456@127.0.0.1:5000/v2/service_instances/123 -d '{"organization_guid": "org-guid-here", "plan_id": "C8E570F3-17E5-4F69-A4C3-6789E6FA62E8","service_id":"45DDB6EA-7EB1-46B5-8E9A-514A7078A6EA"}'
sleep 3
curl -XPUT  michael:123456@127.0.0.1:5000/v2/service_instances/123/service_bindings/456 -d '{"app_gui": "org-guid-here", "plan_id": "C8E570F3-17E5-4F69-A4C3-6789E6FA62E8","service_id":"45DDB6EA-7EB1-46B5-8E9A-514A7078A6EA"}'
#sleep 3
curl -XDELETE michael:123456@127.0.0.1:5000/v2/service_instances/123?service_id=45DDB6EA-7EB1-46B5-8E9A-514A7078A6EA\&plan_id=C8E570F3-17E5-4F69-A4C3-6789E6FA62E8
#curl -XDELETE -v michael:123456@127.0.0.1:5000/v2/service_instances/123/service_bindings/456