import 'package:flutter/material.dart';
import "package:namer_app/globalData.dart";

import "../custom_widget.dart";
import 'exploreDisplayActionsOrReactions.dart';

class FillServicesCard {
  static List<Widget> serviceCardList = [];
  static int serviceIndex = 0;

  static Widget serviceTiles(int index) {
    if (index + 1 < Globaldata.serviceList.length && index % 2 == 0) {
      serviceIndex++;
      return Row(
        children: [
          Padding(padding: EdgeInsets.only(right: 5)),
          SquareServicesCard(
            name: Globaldata.serviceList[index].name,
            color: HexColor(Globaldata.serviceList[index].color),
            imageUrl: Globaldata.serviceList[index].icon,
            onTap: () {
              Navigator.push(Globaldata.myContext, MaterialPageRoute(builder: (context) => DisplayActionOrReactionPage(serviceIndex: index,)));
            },
          ),
          SquareServicesCard(
            name: Globaldata.serviceList[index + 1].name,
            color: HexColor(Globaldata.serviceList[index + 1].color),
            imageUrl: Globaldata.serviceList[index + 1].icon,
            onTap: () {
              Navigator.push(Globaldata.myContext, MaterialPageRoute(builder: (context) => DisplayActionOrReactionPage(serviceIndex: index + 1,)));
            },
          ),
        ],
      );
    }
    return Row(
      children: [
        Padding(padding: EdgeInsets.only(right: 5)),
        SquareServicesCard(
          name: Globaldata.serviceList[index].name,
          color: HexColor(Globaldata.serviceList[index].color),
          imageUrl: Globaldata.serviceList[index].icon,
          onTap: () {
            Navigator.push(Globaldata.myContext, MaterialPageRoute(builder: (context) => DisplayActionOrReactionPage(serviceIndex: index,)));
          },
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
