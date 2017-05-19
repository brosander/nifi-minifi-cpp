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


#include <cassert>
#include <chrono>
#include <fstream>
#include <memory>
#include <string>
#include <thread>
#include <type_traits>
#include <vector>
#include <sys/stat.h>
#include <uuid/uuid.h>
#include <fstream>
#include "FlowController.h"
#include "unit/ProvenanceTestHelper.h"
#include "processors/GetFile.h"
#include "core/Core.h"
#include "core/FlowFile.h"
#include "core/Processor.h"
#include "core/ProcessContext.h"
#include "core/ProcessSession.h"
#include "core/ProcessorNode.h"
#include "TestBase.h"

int main(int argc, char  **argv)
{
  TestController testController;
  std::shared_ptr<core::Processor> processor = std::make_shared<
      org::apache::nifi::minifi::processors::ExecuteProcess>("executeProcess");
  processor->setMaxConcurrentTasks(1);

  std::shared_ptr<core::Repository> test_repo =
      std::make_shared<TestRepository>();

  std::shared_ptr<TestRepository> repo =
      std::static_pointer_cast<TestRepository>(test_repo);
  std::shared_ptr<minifi::FlowController> controller = std::make_shared<
      TestFlowController>(test_repo, test_repo);

  uuid_t processoruuid;
  assert(true == processor->getUUID(processoruuid));

  std::shared_ptr<minifi::Connection> connection = std::make_shared<
      minifi::Connection>(test_repo, "executeProcessConnection");
  connection->setRelationship(core::Relationship("success", "description"));

  // link the connections so that we can test results at the end for this
  connection->setSource(processor);
  connection->setDestination(processor);

  connection->setSourceUUID(processoruuid);
  connection->setDestinationUUID(processoruuid);

  processor->addConnection(connection);
  assert(processor->getName() == "executeProcess");

  std::shared_ptr<core::FlowFile> record;
  processor->setScheduledState(core::ScheduledState::RUNNING);

  processor->initialize();

  std::atomic<bool> is_ready(false);

  std::vector<std::thread> processor_workers;

  core::ProcessorNode node2(processor);
  std::shared_ptr<core::ProcessContext> contextset = std::make_shared<
      core::ProcessContext>(node2, test_repo);
  core::ProcessSessionFactory factory(contextset.get());
  processor->onSchedule(contextset.get(), &factory);

  for (int i = 0; i < 1; i++) {
    //
    processor_workers.push_back(
        std::thread(
            [processor,test_repo,&is_ready]()
            {
              core::ProcessorNode node(processor);
              std::shared_ptr<core::ProcessContext> context = std::make_shared<core::ProcessContext>(node, test_repo);
              context->setProperty(org::apache::nifi::minifi::processors::ExecuteProcess::Command,"sleep 0.5");
              //context->setProperty(org::apache::nifi::minifi::processors::ExecuteProcess::CommandArguments," 1 >>" + ss.str());
              std::shared_ptr<core::ProcessSession> session = std::make_shared<core::ProcessSession>(context.get());
              while(!is_ready.load(std::memory_order_relaxed)) {

              }

              processor->onTrigger(context.get(), session.get());

            }));
  }

  is_ready.store(true, std::memory_order_relaxed);
  //is_ready.store(true);

  std::for_each(processor_workers.begin(), processor_workers.end(),
                [](std::thread &t)
                {
                  t.join();
                });


  std::shared_ptr<org::apache::nifi::minifi::processors::ExecuteProcess> execp =
      std::static_pointer_cast<
          org::apache::nifi::minifi::processors::ExecuteProcess>(processor);

}
