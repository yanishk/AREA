import 'package:flutter/material.dart';
import "package:namer_app/globalData.dart";
import 'dart:developer' as developer;

import "exploreWidget.dart";
import "../custom_widget.dart";

class FillServicesAndWorkflow {
  static List<Widget> servicesAndWorkflow = [];

  static Widget serviceAndWorkflowTiles(int index, int serviceIndex) {
    if (index % 2 == 0) {
      if (serviceIndex < Globaldata.serviceList.length) {
        developer.log("Globaldata.serviceList[serviceIndex].icon : ${Globaldata.serviceList[serviceIndex].icon}");
        return ServicesCard(
          name: Globaldata.serviceList[serviceIndex].name,
          color: HexColor(Globaldata.serviceList[serviceIndex].color),
          imageUrl: Globaldata.serviceList[serviceIndex].icon,
        );
      }
    }
    return WorkflowCard();
  }

  static void fillServicesAndWorkflowList() {
    servicesAndWorkflow = [];
    int serviceIndex = 0;
    for (int allIndex = 0; allIndex < Globaldata.serviceList.length * 2; allIndex++) {
      servicesAndWorkflow.add(serviceAndWorkflowTiles(allIndex, serviceIndex));
      if (allIndex % 2 == 0) {
        if (serviceIndex < Globaldata.serviceList.length) {
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
      name: Globaldata.serviceList[index].name,
      color: HexColor(Globaldata.serviceList[index].color),
      imageUrl: Globaldata.serviceList[index].icon,
    );
  }

  static void fillServicesCardList () {
    serviceCardList = [];
    for (int i = 0; i < Globaldata.serviceList.length; i++) {
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