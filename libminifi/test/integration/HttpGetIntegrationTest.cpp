/**
 *
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <sys/stat.h>
#include <cassert>
#include <utility>
#include <chrono>
#include <fstream>
#include <memory>
#include <string>
#include <thread>
#include <type_traits>
#include <vector>
#include "../TestBase.h"
#include "utils/StringUtils.h"
#include "core/Core.h"
#include "../include/core/logging/Logger.h"
#include "core/ProcessGroup.h"
#include "core/yaml/YamlConfiguration.h"
#include "FlowController.h"
#include "properties/Configure.h"
#include "../unit/ProvenanceTestHelper.h"
#include "io/StreamFactory.h"

void waitToVerifyProcessor() {
  std::this_thread::sleep_for(std::chrono::seconds(10));
}

int main(int argc, char **argv) {
  LogTestController::getInstance().setInfo<minifi::processors::InvokeHTTP>();
  LogTestController::getInstance().setInfo<minifi::processors::LogAttribute>();
  std::string key_dir, test_file_location;
  if (argc > 1) {
    test_file_location = argv[1];
    key_dir = argv[2];
  }
  std::shared_ptr<minifi::Configure> configuration = std::make_shared<minifi::Configure>();
  configuration->set(minifi::Configure::nifi_default_directory, key_dir);
  mkdir("content_repository", S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH);

  std::shared_ptr<core::Repository> test_repo = std::make_shared<TestRepository>();
  std::shared_ptr<core::Repository> test_flow_repo = std::make_shared<TestFlowRepository>();

  configuration->set(minifi::Configure::nifi_flow_configuration_file, test_file_location);

  std::shared_ptr<minifi::io::StreamFactory> stream_factory = std::make_shared<minifi::io::StreamFactory>(configuration);
  std::unique_ptr<core::FlowConfiguration> yaml_ptr = std::unique_ptr<core::YamlConfiguration>(new core::YamlConfiguration(test_repo, test_repo, stream_factory, configuration, test_file_location));
  std::shared_ptr<TestRepository> repo = std::static_pointer_cast<TestRepository>(test_repo);

  std::shared_ptr<minifi::FlowController> controller = std::make_shared<minifi::FlowController>(test_repo, test_flow_repo, configuration, std::move(yaml_ptr), DEFAULT_ROOT_GROUP_NAME,
  true);

  core::YamlConfiguration yaml_config(test_repo, test_repo, stream_factory, configuration, test_file_location);

  std::unique_ptr<core::ProcessGroup> ptr = yaml_config.getRoot(test_file_location);
  std::shared_ptr<core::ProcessGroup> pg = std::shared_ptr<core::ProcessGroup>(ptr.get());
  ptr.release();

  controller->load();
  controller->start();
  waitToVerifyProcessor();

  controller->waitUnload(60000);
  std::string logs = LogTestController::getInstance().log_output.str();
  assert(logs.find("key:filename value:") != std::string::npos);
  assert(logs.find("key:invokehttp.request.url value:https://raw.githubusercontent.com/curl/curl/master/docs/examples/httpput.c") != std::string::npos);
  assert(logs.find("Size:3734 Offset:0") != std::string::npos);
  assert(logs.find("key:invokehttp.status.code value:200") != std::string::npos);
  std::string stringtofind = "Resource Claim created ./content_repository/";

  size_t loc = logs.find(stringtofind);
  while (loc > 0) {
    std::string id = logs.substr(loc + stringtofind.size(), 36);

    loc = logs.find(stringtofind, loc + 1);
    std::string path = "content_repository/" + id;
    unlink(path.c_str());

    if (loc == std::string::npos)
      break;
  }
  LogTestController::getInstance().reset();
  rmdir("./content_repository");
  return 0;
}
