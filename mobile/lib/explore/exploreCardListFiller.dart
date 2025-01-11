import 'package:flutter/material.dart';
import "package:namer_app/globalData.dart";

import "exploreWidget.dart";
import "../custom_widget.dart";

class FillServicesCard {
  static List<Widget> serviceCardList = [];
  static int serviceIndex = 0;

  static Widget serviceTiles(int index) {
    if (index + 1 < Globaldata.serviceList.length && index % 2 == 0) {
      serviceIndex++;
      return Row(
        children: [
          MySquareServicesCard(
            name: Globaldata.serviceList[index].name,
            color: HexColor(Globaldata.serviceList[index].color),
            imageUrl: Globaldata.serviceList[index].icon,
            index: index++,
          ),
          MySquareServicesCard(
            name: Globaldata.serviceList[index].name,
            color: HexColor(Globaldata.serviceList[index].color),
            imageUrl: Globaldata.serviceList[index].icon,
            index: index,
          ),
        ],
      );
    }
    return Row(
      children: [
        MySquareServicesCard(
          name: Globaldata.serviceList[index].name,
          color: HexColor(Globaldata.serviceList[index].color),
          imageUrl: Globaldata.serviceList[index].icon,
        ),
      ],
    );
  }

  static void fillServicesCardList () {
    serviceCardList = [];
    for (serviceIndex = 0; serviceIndex < Globaldata.serviceList.length; serviceIndex++) {
      serviceCardList.add(serviceTiles(serviceIndex));
    }
  }
}
