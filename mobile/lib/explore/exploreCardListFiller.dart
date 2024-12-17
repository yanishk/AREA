import 'package:flutter/material.dart';

import "exploreWidget.dart";
import "../custom_widget.dart";
import "explorePageData.dart";

class FillServicesAndWorkflow {
  static List<Widget> servicesAndWorkflow = [];

  static Widget serviceAndWorkflowTiles(int index, int serviceIndex) {
    if (index % 2 == 0) {
      if (serviceIndex < AllServicesList.allServiceNameList.length) {
        return ServicesCard(
          name: AllServicesList.allServiceNameList[serviceIndex],
          color: HexColor(AllServicesList.allServiceColorList[serviceIndex]),
          imageUrl: AllServicesList.allServiceIconsList[serviceIndex],
        );
      }
    }
    return WorkflowCard();
  }

  static void fillServicesAndWorkflowList() {
    servicesAndWorkflow = [];
    int serviceIndex = 0;
    for (int allIndex = 0; allIndex < 5; allIndex++) {
      servicesAndWorkflow.add(serviceAndWorkflowTiles(allIndex, serviceIndex));
      if (allIndex % 2 == 0) {
        if (serviceIndex < AllServicesList.allServiceNameList.length) {
          serviceIndex++;
        }
      }
    }
  }
}

class FillServicesCard {
  static List<Widget> serviceCardList = [];

  static Widget serviceTiles(int index) {
    return ServicesCard(
      name: AllServicesList.allServiceNameList[index],
      color: HexColor(AllServicesList.allServiceColorList[index]),
      imageUrl: AllServicesList.allServiceIconsList[index],
    );
  }

  static void fillServicesCardList () {
    serviceCardList = [];
    for (int i = 0; i < AllServicesList.allServiceNameList.length; i++) {
      serviceCardList.add(serviceTiles(i));
    }
  }
}

class FillWorkflowsCard {
  static List<Widget> workflowsCardList = [];

  static Widget workflowsTiles(int index) {
    return WorkflowCard();
  }

  static void fillWorkflowsCardList () {
    workflowsCardList = [];
    for (int i = 0; i < 5; i++) {
      workflowsCardList.add(workflowsTiles(i));
    }
  }
}